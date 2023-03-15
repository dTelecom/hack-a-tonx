package main

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"github.com/carlmjohnson/requests"
	"github.com/lucsky/cuid"
	"github.com/rs/zerolog/log"
	"main/pkg/node"
	"main/pkg/relay"
	"main/pkg/sfu"
	"main/pkg/ton"
	"math"
	"strconv"
	"sync"
	"time"
)

// Rooms global map rooms
var Rooms sync.Map

// Handlers global pubsub handlers
var Handlers sync.Map

// Room for participants
type Room struct {
	SID                 string
	OnlineParticipants  sync.Map
	EndedParticipants   sync.Map
	Session             sfu.Session
	Hosts               sync.Map
	Node                node.Node
	RemoteViewersCount  sync.Map
	LocalViewersCount   int
	ClientAddress       string
	ClientPk            ed25519.PublicKey
	URL                 string
	CallID              string
	FirstNotifyResponse NotifyResponse
	LastNotifyResponse  NotifyResponse
	opened              bool
	created             bool
	closed              bool
	ended               bool
	createdChan         chan struct{}
}

// RoomMessage typed json from participant
type RoomMessage struct {
	Participant *Participant    `json:"participant"`
	Payload     json.RawMessage `json:"payload"`
}

// RelayMessage data
type RelayMessage struct {
	Data   []byte
	PeerID string
	ID     string
	Host   string
}

// NotifyData data
type NotifyData struct {
	Duration int    `json:"duration"`
	SID      string `json:"sid"`
	CallID   string `json:"callID"`
	UID      string `json:"uid"`
	Type     string `json:"type"`
}

// NotifyRequest data
type NotifyRequest struct {
	Message   []byte `json:"message"`
	Signature []byte `json:"signature"`
}

// NotifyResponse data
type NotifyResponse struct {
	Message   []byte `json:"message"`
	Signature []byte `json:"signature"`
	Duration  int    `json:"duration"`
}

// ParticipantsMessage between nodes
type ParticipantsMessage struct {
	Participants map[string]*Participant `json:"participants"`
	ViewersCount int                     `json:"viewersCount"`
}

// ParticipantsCount all
type ParticipantsCount struct {
	ParticipantsCount int `json:"participantsCount"`
	ViewersCount      int `json:"viewersCount"`
}

// OnJoin participant
func (r *Room) OnJoin(participant *Participant) {
	if participant.NoPublish == false {
		r.Broadcast(participant, "onJoin", nil)
	} else {
		r.LocalViewersCount++
	}
	r.BroadcastLocal("participantsCount", r.GetAllCountJson())
	go r.NotifyAndTx(participant, "join")
}

// OnJoinRemote participant
func (r *Room) OnJoinRemote(participant *Participant) {
	if participant.NoPublish == false {
		roomMessage := &RoomMessage{Participant: participant, Payload: nil}
		r.BroadcastLocal("onJoin", roomMessage)
		r.BroadcastLocal("participantsCount", r.GetAllCountJson())
	}
}

// OnLeave participant
func (r *Room) OnLeave(participant *Participant) {
	if participant.NoPublish == false {
		r.Broadcast(participant, "onLeave", nil)
		r.EndedParticipants.Store(participant, true)
	} else {
		if r.LocalViewersCount > 0 {
			r.LocalViewersCount--
		}
	}
	r.BroadcastLocal("participantsCount", r.GetAllCountJson())
	go r.NotifyAndTx(participant, "leave")
}

// OnLeaveRemote participant
func (r *Room) OnLeaveRemote(participant *Participant) {
	if participant.NoPublish == false {
		roomMessage := &RoomMessage{Participant: participant, Payload: nil}
		r.BroadcastLocal("onLeave", roomMessage)
		r.BroadcastLocal("participantsCount", r.GetAllCountJson())
	}
}

// OnStream participant
func (r *Room) OnStream(participant *Participant) {
	if participant.NoPublish == false {
		r.Broadcast(participant, "onStream", nil)
	}
}

// OnEnd participant
func (r *Room) OnEnd(participant *Participant) {
	if participant.IsHost {
		r.Broadcast(participant, "end", nil)
		r.EndRoom()
	}
}

// OnRemoteEnd participant
func (r *Room) OnRemoteEnd(participant *Participant) {
	if participant.IsHost {
		r.EndRoom()
	}
}

// EndRoom participant
func (r *Room) EndRoom() {
	r.ended = true
	r.OnlineParticipants.Range(func(_, ival interface{}) bool {
		participant, _ := ival.(*Participant)
		if participant.Host == r.Node.ID().Pretty() {
			participant.Close()
		}
		return true
	})
}

// Broadcast all room participants
func (r *Room) Broadcast(participant *Participant, method string, params json.RawMessage) {
	roomMessage := &RoomMessage{Participant: participant, Payload: params}
	r.BroadcastLocal(method, roomMessage)
	r.Publish(method, roomMessage)
}

// BroadcastLocal all local participants
func (r *Room) BroadcastLocal(method string, params interface{}) {
	r.OnlineParticipants.Range(func(_, ival interface{}) bool {
		participant, _ := ival.(*Participant)
		if participant.Host == r.Node.ID().Pretty() {
			if err := participant.conn.Notify(participant.ctx, method, params); err != nil {
				log.Error().Err(err).Msg(method)
			}
		}
		return true
	})
}

// OnRemoteMessage from p2p
func (r *Room) OnRemoteMessage(senderID string, pubMessage *node.PubMessage) {

	if senderID == r.Node.ID().Pretty() {
		return
	}

	switch pubMessage.Method {
	case "internal":
		var participantsMessage ParticipantsMessage
		err := json.Unmarshal(pubMessage.Payload, &participantsMessage)
		if err != nil {
			log.Printf("OnRemoteMessage err: %v %v", err, pubMessage)
			return
		}
		r.OnRemoteParticipants(senderID, participantsMessage.Participants)
		r.OnRemoteViewers(senderID, participantsMessage.ViewersCount)
	case "relayOffer":
		var relayOffer RelayMessage
		err := json.Unmarshal(pubMessage.Payload, &relayOffer)
		if err != nil {
			log.Error().Err(err).Msg("relayOffer")
			return
		}

		if relayOffer.Host != r.Node.ID().Pretty() {
			return
		}

		data, derr := r.Session.AddRelayPeer(relayOffer.PeerID, relayOffer.Data)
		log.Printf("relayAnswer: %v %v", string(data), derr)

		var relayAnswer RelayMessage
		relayAnswer.Data = data
		relayAnswer.Host = relayOffer.Host
		relayAnswer.ID = relayOffer.ID
		relayAnswer.PeerID = relayOffer.PeerID

		r.Publish("relayAnswer", relayAnswer)
	case "relayAnswer":
		var relayAnswer RelayMessage
		err := json.Unmarshal(pubMessage.Payload, &relayAnswer)
		if err != nil {
			log.Error().Err(err).Msg("relayAnswer")
			return
		}

		log.Printf("relayAnswer: %v %v %v %v", relayAnswer.ID, relayAnswer.Host, relayAnswer.PeerID, string(relayAnswer.Data))

		if ival, ok := Handlers.Load(relayAnswer.ID); ok {
			messages, _ := ival.(chan []byte)
			messages <- relayAnswer.Data
		}
	case "end":
		log.Printf("end: %v", string(pubMessage.Payload))

		var roomMessage RoomMessage
		err := json.Unmarshal(pubMessage.Payload, &roomMessage)
		if err != nil {
			log.Error().Err(err).Msg("end")
			return
		}
		r.OnRemoteEnd(roomMessage.Participant)

	default:
		r.BroadcastLocal(pubMessage.Method, pubMessage.Payload)
	}
}

// OnRemoteParticipants from p2p
func (r *Room) OnRemoteParticipants(hostID string, participants map[string]*Participant) {

	r.Hosts.Store(hostID, true)

	for UID, participant := range participants {
		ival, ok := r.OnlineParticipants.Load(UID)
		if !ok {
			r.OnlineParticipants.Store(UID, participant)
			r.OnJoinRemote(participant)
		} else {
			remoteParticipant := ival.(*Participant)
			if remoteParticipant.StreamID != participant.StreamID {
				remoteParticipant.StreamID = participant.StreamID
			}
		}
	}

	r.OnlineParticipants.Range(func(_, ival interface{}) bool {
		participant, _ := ival.(*Participant)
		if participant.Host == hostID {
			_, ok := participants[participant.UID]
			if !ok {
				r.OnlineParticipants.Delete(participant.UID)
				r.OnLeaveRemote(participant)
			}
		}
		return true
	})
}

// OnRemoteViewers from p2p
func (r *Room) OnRemoteViewers(hostID string, viewersCount int) {

	prevCount := 0
	ival, ok := r.RemoteViewersCount.Load(hostID)
	if ok {
		prevCount = ival.(int)
	}

	if prevCount != viewersCount {
		r.RemoteViewersCount.Store(hostID, viewersCount)
		r.BroadcastLocal("participantsCount", r.GetAllCountJson())
	}
}

// Close room
func (r *Room) Close() {

	if r.closed {
		return
	}
	r.closed = true
	if r.created {
		<-r.createdChan
		duration := r.getEndedDuration()
		log.Printf("duration calc: %v", duration)
		log.Printf("duration last: %v", r.LastNotifyResponse.Duration)

		log.Printf("end call: %v", r.LastNotifyResponse)

		err := ton.EndCall(r.ClientAddress, r.LastNotifyResponse.Signature, r.LastNotifyResponse.Message)
		log.Printf("err: %v", err)
	}
}

// GetParticipants all partiipants
func (r *Room) GetParticipants() map[string]*Participant {
	participants := make(map[string]*Participant)
	r.OnlineParticipants.Range(func(_, ival interface{}) bool {
		participant, _ := ival.(*Participant)
		participants[participant.UID] = participant
		return true
	})
	return participants
}

// GetPublishParticipants all partiipants
func (r *Room) GetPublishParticipants() map[string]*Participant {
	participants := make(map[string]*Participant)
	r.OnlineParticipants.Range(func(_, ival interface{}) bool {
		participant, _ := ival.(*Participant)
		if participant.NoPublish == false {
			participants[participant.UID] = participant
		}
		return true
	})
	return participants
}

// GetLocalParticipants all local publishers
func (r *Room) GetLocalParticipants() map[string]*Participant {
	participants := make(map[string]*Participant)
	r.OnlineParticipants.Range(func(_, ival interface{}) bool {
		participant, _ := ival.(*Participant)
		if participant.Host == r.Node.ID().Pretty() {
			if participant.NoPublish == false {
				participants[participant.UID] = participant
			}
		}
		return true
	})
	return participants
}

// GetLocalViewersCount all local viewers
func (r *Room) GetLocalViewersCount() int {
	count := 0
	r.OnlineParticipants.Range(func(_, ival interface{}) bool {
		participant, _ := ival.(*Participant)
		if participant.Host == r.Node.ID().Pretty() {
			if participant.NoPublish == true {
				count++
			}
		}
		return true
	})
	return count
}

// GetAllCountJson json
func (r *Room) GetAllCountJson() *RoomMessage {
	viewersCount := r.LocalViewersCount
	r.RemoteViewersCount.Range(func(_, ival interface{}) bool {
		viewersCount += ival.(int)
		return true
	})

	publishParticipants := r.GetPublishParticipants()

	all := &ParticipantsCount{
		ParticipantsCount: len(publishParticipants),
		ViewersCount:      viewersCount,
	}
	payload, _ := json.Marshal(all)

	roomMessage := &RoomMessage{Participant: nil, Payload: payload}
	return roomMessage
}

// RelayAll all partiipants
func (r *Room) RelayAll() {
	hosts := make(map[string]bool)
	r.Hosts.Range(func(ikey, _ interface{}) bool {
		host := ikey.(string)
		hosts[host] = true
		return true
	})

	if len(hosts) == 0 {
		return
	}

	participants := r.GetLocalParticipants()
	for host := range hosts {
		signalFunc := func(meta relay.PeerMeta, signal []byte) ([]byte, error) {
			log.Printf("signalFunc: %v %v", meta, string(signal))

			id := cuid.New()
			relayOffer := &RelayMessage{
				Data:   signal,
				PeerID: meta.PeerID,
				ID:     id,
				Host:   host,
			}

			messages := make(chan []byte)
			Handlers.Store(id, messages)

			log.Printf("signalFunc offer: %v %v %v %v", relayOffer.ID, relayOffer.Host, relayOffer.PeerID, string(relayOffer.Data))
			r.Publish("relayOffer", relayOffer)
			message := <-messages
			log.Printf("signalFunc answer: %v", string(message))

			return message, nil
		}

		for _, participant := range participants {
			if participant.StreamID != "" {
				_, ok := participant.relayed[host]
				if ok == false {
					log.Printf("start relay: %v", participant.UID)
					data, derr := participant.Peer.Publisher().Relay(signalFunc)
					log.Printf("relay: %v %v", data, derr)
					participant.relayed[host] = true
				}
			}
		}
	}
}

// Publish pub sub
func (r *Room) Publish(method string, data interface{}) {
	payload, _ := json.Marshal(data)
	pubMessage := &node.PubMessage{Method: method, Payload: payload}
	json, err := json.Marshal(pubMessage)
	if err != nil {
		log.Error().Err(err).Msg("Publish")
	}

	err = r.Node.SendMessage(context.Background(), r.SID, json)
	if err != nil {
		log.Error().Err(err).Msg("Publish")
	}
}

func (r *Room) IsClosed() bool {
	if r.closed {
		return r.closed
	}
	return r.ended
}

func (r *Room) observer() {
	counter := 0

loop:
	for {
		time.Sleep(time.Duration(1) * time.Second)

		participantsMessage := &ParticipantsMessage{
			Participants: r.GetLocalParticipants(),
			ViewersCount: r.GetLocalViewersCount(),
		}

		r.Publish("internal", participantsMessage)

		participants := r.GetParticipants()

		if len(participants) == 0 {
			counter++
		} else {
			counter = 0
		}

		if counter > 5 {
			r.Node.LeaveRoom(r.SID)
			r.Close()
			break loop
		}
		r.RelayAll()
		if r.created == false {
			go r.createCall()
		}
	}
}

func (r *Room) createCall() {
	duration := r.getLocalDuration()
	if duration > 3 {
		r.created = true

		log.Printf("create call: %v", r.FirstNotifyResponse)

		err := ton.CreateCall(r.ClientAddress, r.FirstNotifyResponse.Signature, r.FirstNotifyResponse.Message)
		close(r.createdChan)
		if err != nil {
			log.Printf("err: %v", err)
			r.EndRoom()
		}
	}
}

func (r *Room) getEndedDuration() int {
	duration := 0.0

	r.EndedParticipants.Range(func(k, _ interface{}) bool {
		participant := k.(*Participant)
		difference := participant.RemovedAt.Sub(participant.AddedAt)
		duration += difference.Seconds()
		return true
	})

	minutes := int(math.Ceil(duration / 60.0))
	return minutes
}

func (r *Room) getLocalDuration() int {
	duration := 0.0
	now := time.Now()
	r.OnlineParticipants.Range(func(_, v interface{}) bool {
		participant := v.(*Participant)
		if participant.Host == r.Node.ID().Pretty() {
			difference := now.Sub(participant.AddedAt)
			duration += difference.Seconds()
		}
		return true
	})

	minutes := int(math.Ceil(duration / 60.0))
	return minutes
}

// NotifyAndTx participant
func (r *Room) NotifyAndTx(participant *Participant, action string) {
	duration := r.getEndedDuration()
	log.Printf("duration: %v", duration)

	notifyData := NotifyData{
		Duration: duration,
		SID:      r.SID,
		CallID:   r.CallID,
		UID:      participant.UID,
		Type:     action,
	}

	j, err := json.Marshal(notifyData)
	if err != nil {
		log.Printf("NotifyAndTx: %v", err)
		return
	}

	sign, err := ton.GetSignature(j)
	if err != nil {
		log.Printf("NotifyAndTx: %v", err)
		return
	}

	notifyRequest := NotifyRequest{
		Message:   j,
		Signature: sign,
	}

	var notifyResponse NotifyResponse

	err = requests.
		URL(r.URL).
		BodyJSON(&notifyRequest).
		ToJSON(&notifyResponse).
		Fetch(context.Background())
	if err != nil {
		log.Printf("err: %v", err)
		return
	}

	message := notifyData.CallID + ":" + strconv.Itoa(notifyData.Duration)
	verified := ton.VerifyMessage(r.ClientPk, []byte(message), notifyResponse.Signature)
	if verified != true {
		log.Printf("not verified signature")
		r.EndRoom()
		return
	}

	if action == "join" {
		if r.opened {
			return
		}
		r.opened = true
		r.FirstNotifyResponse = notifyResponse
	}
	if action == "leave" {
		r.LastNotifyResponse = notifyResponse
	}
}

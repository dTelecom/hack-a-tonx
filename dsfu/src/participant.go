package main

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pion/webrtc/v3"
	"github.com/rs/zerolog/log"
	"github.com/sourcegraph/jsonrpc2"
	"main/pkg/node"
	"main/pkg/sfu"
	"main/pkg/ton"
	"time"
)

// Participant participant
type Participant struct {
	Peer *sfu.PeerLocal `json:"-"`
	Node node.Node      `json:"-"`

	SID        string `json:"sid"`
	UID        string `json:"uid"`
	Name       string `json:"name"`
	StreamID   string `json:"streamID"`
	IsHost     bool   `json:"isHost"`
	Host       string `json:"host"`
	NoPublish  bool   `json:"noPublish"`
	AudioMuted bool   `json:"audioMuted"`
	VideoMuted bool   `json:"videoMuted"`

	AddedAt   time.Time `json:"-"`
	RemovedAt time.Time `json:"-"`

	ctx     context.Context `json:"-"`
	conn    *jsonrpc2.Conn  `json:"-"`
	relayed map[string]bool `json:"-"`
}

// NewParticipant create new JSONSignal
func NewParticipant(peer *sfu.PeerLocal, node node.Node) *Participant {
	return &Participant{
		Peer:    peer,
		Node:    node,
		relayed: make(map[string]bool),
	}
}

// JoinRequest message sent when initializing a peer connection
type JoinRequest struct {
	Token     string                    `json:"token"`
	Signature string                    `json:"signature"`
	Name      string                    `json:"name"`
	Offer     webrtc.SessionDescription `json:"offer"`
}

// Token model
type Token struct {
	SID           string `json:"sid"`
	UID           string `json:"uid"`
	Name          string `json:"name"`
	IsHost        bool   `json:"isHost"`
	ClientAddress string `json:"clientAddress"`
	URL           string `json:"url"`
	CallID        string `json:"callID"`
	NoPublish     bool   `json:"noPublish"`
}

// WebrtcNegotiation message sent when renegotiating the peer connection
type WebrtcNegotiation struct {
	Desc webrtc.SessionDescription `json:"desc"`
}

// WebrtcTrickle message sent when renegotiating the peer connection
type WebrtcTrickle struct {
	Target    int                     `json:"target"`
	Candidate webrtc.ICECandidateInit `json:"candidate"`
}

// MuteEvent model
type MuteEvent struct {
	Kind  string `json:"kind"`
	Muted bool   `json:"muted"`
}

// Handle incoming RPC call events like join, answer, offer and trickle
func (p *Participant) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	replyError := func(err error) {
		log.Error().Err(err).Msg("replyError")
		_ = conn.ReplyWithError(ctx, req.ID, &jsonrpc2.Error{
			Code:    500,
			Message: fmt.Sprintf("%s", err),
		})
	}

	log.Info().Str("message", req.Method).Msg("received")

	switch req.Method {
	case "join":
		if p.UID != "" {
			err := fmt.Errorf("already joined")
			replyError(err)
			break
		}

		var joinRequest JoinRequest
		err := json.Unmarshal(*req.Params, &joinRequest)
		if err != nil {
			replyError(err)
			break
		}

		log.Printf("got joinRequest: %v", joinRequest)

		tokenJson, err := base64.StdEncoding.DecodeString(joinRequest.Token)
		if err != nil {
			replyError(err)
			break
		}

		signature, err := base64.StdEncoding.DecodeString(joinRequest.Signature)
		if err != nil {
			replyError(err)
			break
		}

		var token Token
		err = json.Unmarshal(tokenJson, &token)
		if err != nil {
			replyError(err)
			break
		}
		log.Printf("got token: %v", token)

		var room *Room
		if ival, ok := Rooms.Load(token.SID); ok {
			room, _ = ival.(*Room)
			if room.IsClosed() {
				replyError(fmt.Errorf("room closed"))
				break
			}
		}

		var clientPk ed25519.PublicKey
		if room != nil {
			clientPk = room.ClientPk
		} else {
			clientPk, err = ton.GetClientPubKey(token.ClientAddress)
			if err != nil {
				replyError(err)
				break
			}
		}

		verified := ton.VerifyMessage(clientPk, tokenJson, signature)
		if verified != true {
			replyError(fmt.Errorf("not verified signature"))
			break
		}

		log.Printf("verified: %v", verified)

		p.UID = token.UID
		p.SID = token.SID
		p.Name = token.Name
		p.IsHost = token.IsHost

		p.ctx = ctx
		p.conn = conn
		p.Host = p.Node.ID().Pretty()
		p.AddedAt = time.Now()
		p.NoPublish = token.NoPublish
		p.AudioMuted = true
		p.VideoMuted = true

		p.Peer.OnOffer = func(offer *webrtc.SessionDescription) {
			if err := conn.Notify(ctx, "offer", offer); err != nil {
				log.Error().Err(err).Msg("join")
			}
		}

		p.Peer.OnIceCandidate = func(candidate *webrtc.ICECandidateInit, target int) {
			if err := conn.Notify(ctx, "trickle", WebrtcTrickle{
				Candidate: *candidate,
				Target:    target,
			}); err != nil {
				log.Error().Err(err).Msg("join")
			}
		}

		joinConfig := sfu.JoinConfig{
			NoPublish:       false,
			NoSubscribe:     false,
			NoAutoSubscribe: false,
		}

		err = p.Peer.Join(p.SID, p.UID, joinConfig)
		if err != nil {
			replyError(err)
			break
		}

		answer, err := p.Peer.Answer(joinRequest.Offer)
		if err != nil {
			replyError(err)
			break
		}

		if err := conn.Reply(ctx, req.ID, answer); err != nil {
			log.Error().Err(err).Msg("join")
		}

		if room == nil {
			room = &Room{
				SID:           token.SID,
				Session:       p.Peer.Session(),
				Node:          p.Node,
				ClientAddress: token.ClientAddress,
				ClientPk:      clientPk,
				URL:           token.URL,
				CallID:        token.CallID,
				createdChan:   make(chan struct{}),
			}
			Rooms.Store(token.SID, room)
			err := p.Node.JoinRoom(token.SID, p.Node.ID().Pretty(), room.OnRemoteMessage)
			if err != nil {
				log.Error().Err(err).Msg("JoinRoom")
			}

			go room.observer()
		}

		onPublisherTrack := func(track sfu.PublisherTrack) {
			if p.StreamID != track.Track.StreamID() {
				p.StreamID = track.Track.StreamID()
				room.OnStream(p)
			}
		}

		p.Peer.Publisher().OnPublisherTrack(onPublisherTrack)

		room.OnlineParticipants.Store(token.UID, p)
		room.OnJoin(p)

		if err := conn.Notify(ctx, "participants", room.GetPublishParticipants()); err != nil {
			log.Error().Err(err).Msg("join")
		}

	case "offer":
		var webrtcNegotiation WebrtcNegotiation
		err := json.Unmarshal(*req.Params, &webrtcNegotiation)
		if err != nil {
			replyError(err)
			p.Close()
			break
		}

		answer, err := p.Peer.Answer(webrtcNegotiation.Desc)
		if err != nil {
			replyError(err)
			p.Close()
			break
		}
		if err := conn.Reply(ctx, req.ID, answer); err != nil {
			log.Error().Err(err).Msg("offer")
		}

	case "answer":
		var webrtcNegotiation WebrtcNegotiation
		err := json.Unmarshal(*req.Params, &webrtcNegotiation)
		if err != nil {
			replyError(err)
			p.Close()
			break
		}

		err = p.Peer.SetRemoteDescription(webrtcNegotiation.Desc)
		if err != nil {
			replyError(err)
			p.Close()
		}

	case "trickle":
		var webrtcTrickle WebrtcTrickle
		err := json.Unmarshal(*req.Params, &webrtcTrickle)
		if err != nil {
			replyError(err)
			p.Close()
			break
		}

		err = p.Peer.Trickle(webrtcTrickle.Candidate, webrtcTrickle.Target)
		if err != nil {
			replyError(err)
		}

	case "muteEvent":
		var muteEvent MuteEvent
		err := json.Unmarshal(*req.Params, &muteEvent)
		if err != nil {
			replyError(err)
			break
		}
		if muteEvent.Kind == "audio" {
			p.AudioMuted = muteEvent.Muted
		}
		if muteEvent.Kind == "video" {
			p.VideoMuted = muteEvent.Muted
		}

		if p.UID == "" {
			err := fmt.Errorf("not joined")
			replyError(err)
			break
		}
		var room *Room
		if ival, ok := Rooms.Load(p.SID); ok {
			room, _ = ival.(*Room)
		} else {
			err := fmt.Errorf("room not found")
			replyError(err)
			break
		}

		room.Broadcast(p, req.Method, *req.Params)
	case "end":
		if p.UID == "" {
			err := fmt.Errorf("not joined")
			replyError(err)
			break
		}
		var room *Room
		if ival, ok := Rooms.Load(p.SID); ok {
			room, _ = ival.(*Room)
		} else {
			err := fmt.Errorf("room not found")
			replyError(err)
			break
		}
		if p.IsHost == false {
			p.Close()
		} else {
			room.OnEnd(p)
		}
	}
}

// Close ws close
func (p *Participant) Close() {
	if ival, ok := Rooms.Load(p.SID); ok {
		room, _ := ival.(*Room)
		room.OnlineParticipants.Delete(p.UID)
		Rooms.Store(p.SID, room)
		p.RemovedAt = time.Now()
		room.OnLeave(p)
	}

	p.Peer.Close()

	log.Info().Msg("close ws")
}

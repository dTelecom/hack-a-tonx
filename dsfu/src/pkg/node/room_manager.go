package node

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
)

// PubMessage typed json from node
type PubMessage struct {
	Method  string          `json:"method"`
	Payload json.RawMessage `json:"payload"`
}

// RoomMessage holds data to be published in a topic.
type RoomMessage struct {
	Payload  json.RawMessage `json:"payload,omitempty"`
	SenderID peer.ID         `json:"sender_id"`
}

// Room holds room event and pubsub data.
type Room struct {
	name         string
	topic        *pubsub.Topic
	subscription *pubsub.Subscription
	onMessage    OnMessage
}

func newRoom(name string, topic *pubsub.Topic, subscription *pubsub.Subscription, onMessage OnMessage) *Room {
	return &Room{
		name:         name,
		topic:        topic,
		subscription: subscription,
		onMessage:    onMessage,
	}
}

// RoomManager manages rooms through pubsub subscription and implements room operations.
type RoomManager struct {
	ps     *pubsub.PubSub
	node   Node
	kadDHT *dht.IpfsDHT
	rooms  map[string]*Room
	lock   sync.RWMutex
}

// NewRoomManager creates a new room manager.
func NewRoomManager(node Node, kadDHT *dht.IpfsDHT, ps *pubsub.PubSub) *RoomManager {
	mngr := &RoomManager{
		ps:     ps,
		node:   node,
		kadDHT: kadDHT,
		rooms:  make(map[string]*Room),
	}

	return mngr
}

// JoinAndSubscribe joins and subscribes to a room.
func (r *RoomManager) JoinAndSubscribe(roomName string, nickname string, onMessage OnMessage) error {
	if r.HasJoined(roomName) {
		return errors.New("already joined")
	}

	topicName := r.TopicName(roomName)
	topic, err := r.ps.Join(topicName)
	if err != nil {
		return err
	}

	subscription, err := topic.Subscribe()
	if err != nil {
		if topic != nil {
			_ = topic.Close()
		}
		if subscription != nil {
			subscription.Cancel()
		}

		return err
	}

	room := newRoom(roomName, topic, subscription, onMessage)

	r.putRoom(room)
	go r.roomSubscriptionHandler(room)

	return nil
}

// Leave room.
func (r *RoomManager) Leave(roomName string) error {
	if !r.HasJoined(roomName) {
		return errors.New("not joined")
	}

	room, _ := r.getRoom(roomName)

	if room.subscription != nil {
		room.subscription.Cancel()
	}

	if room.topic != nil {
		room.topic.Close()
	}

	r.removeRoom(roomName)

	return nil
}

// HasJoined returns whether the manager has joined a given room.
func (r *RoomManager) HasJoined(roomName string) bool {
	r.lock.RLock()
	defer r.lock.RUnlock()

	_, found := r.rooms[r.TopicName(roomName)]
	return found
}

// TopicName builds a string containing the name of the pubsub topic for a given room name.
func (r *RoomManager) TopicName(roomName string) string {
	return fmt.Sprintf("webrtc/room/%s", roomName)
}

// SendMessage sends a message to a given room.
// Fails if it has not yet joined the given room.
func (r *RoomManager) SendMessage(ctx context.Context, roomName string, msg []byte) error {
	room, found := r.getRoom(roomName)
	if !found {
		return errors.New(fmt.Sprintf("must join the room before sending messages"))
	}

	rm := &RoomMessage{
		Payload:  msg,
		SenderID: r.node.ID(),
	}

	if err := r.publishRoomMessage(ctx, room, rm); err != nil {
		return err
	}

	return nil
}

func (r *RoomManager) putRoom(room *Room) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.rooms[room.topic.String()] = room
}

func (r *RoomManager) getRoom(roomName string) (*Room, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	room, found := r.rooms[r.TopicName(roomName)]
	return room, found
}

func (r *RoomManager) removeRoom(roomName string) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	delete(r.rooms, r.TopicName(roomName))
}

func (r *RoomManager) roomSubscriptionHandler(room *Room) {
loop:
	for {
		subMsg, err := room.subscription.Next(context.Background())
		if err != nil {
			break loop
		}

		if subMsg.ReceivedFrom == r.node.ID() {
			continue
		}

		var rm RoomMessage
		err = json.Unmarshal(subMsg.Data, &rm)
		if err != nil {
			continue
		}

		var pubMessage PubMessage
		if err := json.Unmarshal(rm.Payload, &pubMessage); err != nil {
			continue
		}
		go room.onMessage(rm.SenderID.Pretty(), &pubMessage)
	}
}

func (r *RoomManager) publishRoomMessage(
	ctx context.Context,
	room *Room,
	rm *RoomMessage,
) error {
	rmJSON, err := json.Marshal(rm)
	if err != nil {
		return errors.Wrap(err, "marshalling message")
	}

	if err := room.topic.Publish(ctx, rmJSON); err != nil {
		return err
	}

	return nil
}

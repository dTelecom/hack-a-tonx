package node

import (
	"context"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/crypto"
	discovery2 "github.com/libp2p/go-libp2p/core/discovery"
	libp2phost "github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	discoveryNamespace = "/webrtc"
	privKeyFileName    = "libp2p-webrtc.privkey"
)

// OnMessage callback
type OnMessage func(string, *PubMessage)

// Node struct
type Node interface {
	ID() peer.ID

	Start(ctx context.Context, port uint16) error
	Bootstrap(ctx context.Context, nodeAddrs []multiaddr.Multiaddr) error

	SendMessage(ctx context.Context, roomName string, msg []byte) error

	JoinRoom(roomName string, nickname string, onMessage OnMessage) error
	LeaveRoom(roomName string) error
}

type node struct {
	host            libp2phost.Host
	kadDHT          *dht.IpfsDHT
	ps              *pubsub.PubSub
	roomManager     *RoomManager
	privKeyFileName string
}

// NewNode e
func NewNode(privKeyFileName string) Node {
	return &node{
		host:            nil,
		privKeyFileName: privKeyFileName,
	}
}

func (n *node) ID() peer.ID {
	if n.host == nil {
		return ""
	}
	return n.host.ID()
}

func (n *node) Start(ctx context.Context, port uint16) error {

	nodeAddrStrings := []string{fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port)}

	privKey, err := n.getPrivateKey()
	if err != nil {
		return err
	}

	log.Debug().Msg("creating libp2p host")

	host, err := libp2p.New(
		libp2p.ListenAddrStrings(nodeAddrStrings...),
		libp2p.Identity(privKey),
	)
	if err != nil {
		return errors.Wrap(err, "creating libp2p host")
	}
	n.host = host

	ps, err := pubsub.NewGossipSub(ctx, n.host, pubsub.WithMessageSignaturePolicy(pubsub.StrictSign))
	if err != nil {
		return errors.Wrap(err, "creating pubsub")
	}
	n.ps = ps

	p2pAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/p2p/%s", host.ID().Pretty()))
	if err != nil {
		return errors.Wrap(err, "creating host p2p multiaddr")
	}

	var fullAddrs []string
	for _, addr := range host.Addrs() {
		fullAddrs = append(fullAddrs, addr.Encapsulate(p2pAddr).String())
	}

	log.Printf("started node: %v", fullAddrs)
	return nil
}

func (n *node) Bootstrap(ctx context.Context, nodeAddrs []multiaddr.Multiaddr) error {
	var bootstrappers []peer.AddrInfo
	for _, nodeAddr := range nodeAddrs {
		pi, err := peer.AddrInfoFromP2pAddr(nodeAddr)
		if err != nil {
			return errors.Wrap(err, "parsing bootstrapper node address info from p2p address")
		}

		bootstrappers = append(bootstrappers, *pi)
	}

	kadDHT, err := dht.New(
		ctx,
		n.host,
		dht.BootstrapPeers(bootstrappers...),
		dht.ProtocolPrefix(discoveryNamespace),
		dht.Mode(dht.ModeAutoServer),
	)
	if err != nil {
		return errors.Wrap(err, "creating routing DHT")
	}

	n.kadDHT = kadDHT

	if err := kadDHT.Bootstrap(ctx); err != nil {
		return errors.Wrap(err, "bootstrapping DHT")
	}

	// connect to bootstrap nodes, if any
	for _, pi := range bootstrappers {
		if err := n.host.Connect(ctx, pi); err != nil {
			log.Error().Err(err)
		}
	}

	routingDiscovery := drouting.NewRoutingDiscovery(kadDHT)
	dutil.Advertise(ctx, routingDiscovery, discoveryNamespace)

	// try finding more peers
	go func() {
		for {
			peersChan, err := routingDiscovery.FindPeers(
				ctx,
				discoveryNamespace,
				discovery2.Limit(100),
			)
			if err != nil {
				log.Error().Err(err)
				continue
			}

			// read all channel messages to avoid blocking the find peer query
			for range peersChan {
			}

			var peerInfos []string
			for _, peerID := range kadDHT.RoutingTable().ListPeers() {
				peerInfo := n.host.Peerstore().PeerInfo(peerID)
				peerInfos = append(peerInfos, peerInfo.String())
			}

			<-time.After(time.Second * 1)
		}
	}()

	// Setup room manager
	// connect to bootstrap nodes, if any
	roomManager := NewRoomManager(n, n.kadDHT, n.ps)
	n.roomManager = roomManager

	return nil
}

func (n *node) SendMessage(ctx context.Context, roomName string, msg []byte) error {

	if err := n.roomManager.SendMessage(ctx, roomName, msg); err != nil {
		return errors.Wrap(err, "publishing message")
	}

	return nil
}

func (n *node) JoinRoom(roomName string, nickname string, onMessage OnMessage) error {

	if err := n.roomManager.JoinAndSubscribe(roomName, nickname, onMessage); err != nil {
		return err
	}
	return nil
}

func (n *node) LeaveRoom(roomName string) error {

	if err := n.roomManager.Leave(roomName); err != nil {
		return err
	}
	return nil
}

func (n *node) getPrivateKey() (crypto.PrivKey, error) {

	var generate bool

	privKeyBytes, err := ioutil.ReadFile(n.privKeyFileName)
	if os.IsNotExist(err) {
		log.Printf("no identity private key file found.")
		generate = true
	} else if err != nil {
		return nil, err
	}

	if generate {
		privKey, err := n.generateNewPrivKey()
		if err != nil {
			return nil, err
		}

		privKeyBytes, err := crypto.MarshalPrivateKey(privKey)
		if err != nil {
			return nil, errors.Wrap(err, "marshalling identity private key")
		}

		f, err := os.Create(n.privKeyFileName)
		if err != nil {
			return nil, errors.Wrap(err, "creating identity private key file")
		}
		defer f.Close()

		if _, err := f.Write(privKeyBytes); err != nil {
			return nil, errors.Wrap(err, "writing identity private key to file")
		}

		return privKey, nil
	}

	privKey, err := crypto.UnmarshalPrivateKey(privKeyBytes)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling identity private key")
	}

	log.Printf("loaded identity private key from file")

	return privKey, nil
}

func (n *node) generateNewPrivKey() (crypto.PrivKey, error) {

	log.Printf("generating identity private key")
	privKey, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return nil, errors.Wrap(err, "generating identity private key")
	}

	return privKey, nil
}

package ton

import (
	"crypto/ed25519"
	"github.com/rs/zerolog/log"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"math/rand"
	"os"
	"strings"
)

func GetNodeURL() (string, string, ed25519.PublicKey, error) {
	var pk ed25519.PublicKey

	userToncli, err := NewUserToncli(os.Getenv("TON_SEED"), wallet.V4R2, os.Getenv("TON_MASTER_CONTRACT"))
	if err != nil {
		return "", "", pk, err
	}

	nodes, err := userToncli.GetNodeHosts()
	if err != nil {
		return "", "", pk, err
	}

	var keys []string

	for k := range nodes {
		if strings.HasSuffix(k, "tonmeet.com/ws") {
			keys = append(keys, k)
		}
	}

	randomIndex := rand.Intn(len(keys))
	nodeUrl := keys[randomIndex]
	nodeAddress, _ := nodes[nodeUrl]

	log.Printf("keys: %v", keys)
	log.Printf("nodeUrl: %v", nodeUrl)
	log.Printf("nodeAddress: %v", nodeAddress)

	pk, err = userToncli.GetNodePublicKey(nodeAddress)
	if err != nil {
		return "", "", pk, err
	}

	return nodeUrl, nodeAddress.String(), pk, nil
}

func GetSignature(data []byte) ([]byte, error) {
	nodeToncli, err := NewNodeToncli(os.Getenv("TON_SEED"), wallet.V4R2, os.Getenv("TON_MASTER_CONTRACT"))
	if err != nil {
		return nil, err
	}

	return SignMessage(nodeToncli.wallet.PrivateKey(), data), nil
}

func BuildCreateCallMessage(callId uint64) (msg, sign []byte, err error) {
	userToncli, err := NewUserToncli(os.Getenv("TON_SEED"), wallet.V4R2, os.Getenv("TON_MASTER_CONTRACT"))
	if err != nil {
		return nil, nil, err
	}
	return userToncli.BuildCreateCallMessage(callId)
}

func BuildEndCallMessage(callId uint64, spentMinutes uint32) (msg, sign []byte, err error) {
	userToncli, err := NewUserToncli(os.Getenv("TON_SEED"), wallet.V4R2, os.Getenv("TON_MASTER_CONTRACT"))
	if err != nil {
		return nil, nil, err
	}
	return userToncli.BuildEndCallMessage(callId, spentMinutes)
}

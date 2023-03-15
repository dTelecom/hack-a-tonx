package ton

import (
	"crypto/ed25519"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"math/rand"
	"os"
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

	keys := make([]string, len(nodes))

	for k := range nodes {
		keys = append(keys, k)
	}

	randomIndex := rand.Intn(len(keys))
	nodeUrl := keys[randomIndex]
	nodeAddress, _ := nodes[nodeUrl]
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

	return Sign(nodeToncli.wallet.PrivateKey(), data), nil
}

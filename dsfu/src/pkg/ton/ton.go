package ton

import (
	"crypto/ed25519"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"os"
	"strings"
)

func GetClientPubKey(userWalletAddr string) (ed25519.PublicKey, error) {
	var result ed25519.PublicKey
	nodeToncli, err := NewNodeToncli(strings.Split(os.Getenv("TON_SEED"), " "), wallet.V4R2, os.Getenv("TON_MASTER_CONTRACT"))
	if err != nil {
		return result, err
	}

	result, err = nodeToncli.GetUserContractPublicKey(userWalletAddr)
	if err != nil {
		return result, err
	}
	return result, nil
}

func CreateCall(userWalletAddr string, userSign []byte, userMsg []byte) error {
	nodeToncli, err := NewNodeToncli(strings.Split(os.Getenv("TON_SEED"), " "), wallet.V4R2, os.Getenv("TON_MASTER_CONTRACT"))
	if err != nil {
		return err
	}
	return nodeToncli.createCall(userWalletAddr, userSign, userMsg)
}

func GetSignature(data []byte) ([]byte, error) {
	nodeToncli, err := NewNodeToncli(strings.Split(os.Getenv("TON_SEED"), " "), wallet.V4R2, os.Getenv("TON_MASTER_CONTRACT"))
	if err != nil {
		return nil, err
	}

	return Sign(nodeToncli.wallet.PrivateKey(), data), nil
}

func EndCall(userWalletAddr string, userSign []byte, userMsg []byte) error {
	nodeToncli, err := NewNodeToncli(strings.Split(os.Getenv("TON_SEED"), " "), wallet.V4R2, os.Getenv("TON_MASTER_CONTRACT"))
	if err != nil {
		return err
	}
	return nodeToncli.endCall(userWalletAddr, userSign, userMsg)
}

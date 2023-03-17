package ton

import (
	"crypto/ed25519"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"log"
	"os"
)

func GetClientPubKey(userWalletAddr string) (ed25519.PublicKey, error) {
	var result ed25519.PublicKey
	nodeToncli, err := NewNodeToncli(os.Getenv("TON_SEED"), wallet.V4R2, os.Getenv("TON_MASTER_CONTRACT"))
	if err != nil {
		log.Printf("NewNodeToncli %v", os.Getenv("TON_MASTER_CONTRACT"))
		return result, err
	}

	result, err = nodeToncli.GetUserContractPublicKey(userWalletAddr)
	if err != nil {
		log.Printf("GetUserContractPublicKey %v", userWalletAddr)
		return result, err
	}
	return result, nil
}

func CreateCall(userWalletAddr string, userSign []byte, userMsg []byte) error {
	nodeToncli, err := NewNodeToncli(os.Getenv("TON_SEED"), wallet.V4R2, os.Getenv("TON_MASTER_CONTRACT"))
	if err != nil {
		return err
	}
	return nodeToncli.CreateCall(userWalletAddr, userSign, userMsg)
}

func GetSignature(data []byte) ([]byte, error) {
	nodeToncli, err := NewNodeToncli(os.Getenv("TON_SEED"), wallet.V4R2, os.Getenv("TON_MASTER_CONTRACT"))
	if err != nil {
		return nil, err
	}

	return SignMessage(nodeToncli.wallet.PrivateKey(), data), nil
}

func EndCall(userWalletAddr string, userSign []byte, userMsg []byte) error {
	nodeToncli, err := NewNodeToncli(os.Getenv("TON_SEED"), wallet.V4R2, os.Getenv("TON_MASTER_CONTRACT"))
	if err != nil {
		return err
	}
	return nodeToncli.EndCall(userWalletAddr, userSign, userMsg)
}

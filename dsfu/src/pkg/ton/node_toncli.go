package ton

import (
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

type NodeToncli struct {
	api            *ton.APIClient
	wallet         *wallet.Wallet
	masterContract *MasterContract
	contract       *NodeContract
}

func NewNodeToncli(walletSeed []string, walletVersion wallet.Version, masterContractAddr string) (*NodeToncli, error) {
	client := liteclient.NewConnectionPool()

	configUrl := "https://ton-blockchain.github.io/testnet-global.config.json"
	if err := client.AddConnectionsFromConfigUrl(context.Background(), configUrl); err != nil {
		return nil, fmt.Errorf("client.AddConnectionsFromConfigUrl: %w", err)
	}
	api := ton.NewAPIClient(client)

	w, err := wallet.FromSeed(api, walletSeed, walletVersion)
	if err != nil {
		return nil, fmt.Errorf("wallet.FromSeed: %w", err)
	}

	block, err := api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		return nil, fmt.Errorf("api.CurrentMasterchainInfo: %w", err)
	}

	if walletBalance, err := w.GetBalance(context.Background(), block); err != nil {
		return nil, fmt.Errorf("w.GetBalance: %w", err)
	} else {
		fmt.Printf("node wallet (address = %s, balance = %s)\n", w.Address(), walletBalance)
	}

	masterContract := OpenMasterContract(api, address.MustParseAddr(masterContractAddr))
	contractAddr, err := masterContract.GetNodeContractAddress(w.Address())
	if err != nil {
		return nil, err
	}

	fmt.Printf("node contract (address = %s)\n", contractAddr)

	contract := OpenNodeContract(api, contractAddr)
	if nodeData, err := contract.GetData(); err != nil {
		return nil, err
	} else {
		if nodeData.Master.String() != masterContractAddr || nodeData.Owner.String() != w.Address().String() {
			return nil, errors.New("strange node contract data")
		} else {
			if contractBalance, err := contract.GetBalance(); err != nil {
				return nil, fmt.Errorf("contract.GetBalance: %w", err)
			} else {
				fmt.Printf("node contract (address = %s, balance = %s)\n", contractAddr, contractBalance)
			}
			return &NodeToncli{
				api:            api,
				wallet:         w,
				masterContract: masterContract,
				contract:       contract,
			}, nil
		}
	}
}

func (c *NodeToncli) getSignature(data []byte) []byte {
	return Sign(c.wallet.PrivateKey(), data)
}

func (c *NodeToncli) GetUserContractPublicKey(userAddr string) (ed25519.PublicKey, error) {
	userContractAddr, err := c.masterContract.GetUserContractAddress(address.MustParseAddr(userAddr))
	if err != nil {
		return nil, fmt.Errorf("masterContract.GetUserContractAddress: %w", err)
	}
	userContract := OpenUserContract(c.api, userContractAddr)
	userContractData, err := userContract.GetData()
	if err != nil {
		return nil, fmt.Errorf("userContract.GetData: %w", err)
	}
	return userContractData.PublicKey, nil
}

func (c *NodeToncli) createCall(userAddr string, userSign []byte, userMsg []byte) error {
	err := c.contract.SendCreateCall(c.wallet, address.MustParseAddr(userAddr), userSign, userMsg)
	if err != nil {
		err = fmt.Errorf("SendCreateCall: %w", err)
	}
	return err
}

func (c *NodeToncli) endCall(userAddr string, userSign []byte, userMsg []byte) error {
	err := c.contract.SendEndCall(c.wallet, address.MustParseAddr(userAddr), userSign, userMsg)
	if err != nil {
		err = fmt.Errorf("SendEndCall: %w", err)
	}
	return err
}

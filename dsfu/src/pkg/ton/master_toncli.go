package ton

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"strings"
)

type MasterToncli struct {
	api      *ton.APIClient
	contract *MasterContract
}

func NewMasterToncli(contractAddr string) (*MasterToncli, error) {
	client := liteclient.NewConnectionPool()

	configUrl := "https://ton-blockchain.github.io/testnet-global.config.json"
	if err := client.AddConnectionsFromConfigUrl(context.Background(), configUrl); err != nil {
		return nil, fmt.Errorf("client.AddConnectionsFromConfigUrl: %w", err)
	}
	api := ton.NewAPIClient(client)

	contract := OpenMasterContract(api, address.MustParseAddr(contractAddr))
	if hosts, err := contract.GetNodeHosts(); err != nil {
		return nil, fmt.Errorf("masterContract.GetNodeHosts: %w\n", err)
	} else {
		fmt.Printf("hosts = %s\n", hosts)
	}

	return &MasterToncli{
		api:      api,
		contract: contract,
	}, nil
}

func (c *MasterToncli) getWallet(seed string, version wallet.Version) (*wallet.Wallet, error) {
	w, err := wallet.FromSeed(c.api, strings.Split(seed, " "), version)
	if err != nil {
		return nil, fmt.Errorf("wallet.FromSeed: %w", err)
	}

	block, err := c.api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		return nil, fmt.Errorf("api.CurrentMasterchainInfo: %w", err)
	}

	if walletBalance, err := w.GetBalance(context.Background(), block); err != nil {
		return nil, fmt.Errorf("w.GetBalance: %w", err)
	} else {
		fmt.Printf("wallet (address = %s, balance = %s)\n", w.Address(), walletBalance)
	}

	return w, nil
}

func (c *MasterToncli) CreateUser(walletSeed string, walletVersion wallet.Version) error {
	w, err := c.getWallet(walletSeed, walletVersion)
	if err != nil {
		return err
	}
	return c.contract.SendCreateUser(w, w.PrivateKey().Public().(ed25519.PublicKey))
}

func (c *MasterToncli) CreateNode(walletSeed string, walletVersion wallet.Version, nodeHost string) error {
	w, err := c.getWallet(walletSeed, walletVersion)
	if err != nil {
		return err
	}
	return c.contract.SendCreateNode(w, w.PrivateKey().Public().(ed25519.PublicKey), nodeHost)
}

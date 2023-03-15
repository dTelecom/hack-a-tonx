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
	"github.com/xssnick/tonutils-go/tvm/cell"
	"strings"
	"time"
)

type UserToncli struct {
	api            *ton.APIClient
	wallet         *wallet.Wallet
	masterContract *MasterContract
	contract       *UserContract
}

func NewUserToncli(walletSeed string, walletVersion wallet.Version, masterContractAddr string) (*UserToncli, error) {
	client := liteclient.NewConnectionPool()

	configUrl := "https://ton-blockchain.github.io/testnet-global.config.json"
	if err := client.AddConnectionsFromConfigUrl(context.Background(), configUrl); err != nil {
		return nil, fmt.Errorf("client.AddConnectionsFromConfigUrl: %w", err)
	}
	api := ton.NewAPIClient(client)

	w, err := wallet.FromSeed(api, strings.Split(walletSeed, " "), walletVersion)
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
		fmt.Printf("user wallet (address = %s, balance = %s)\n", w.Address(), walletBalance)
	}

	masterContract := OpenMasterContract(api, address.MustParseAddr(masterContractAddr))
	if hosts, err := masterContract.GetNodeHosts(); err != nil {
		return nil, fmt.Errorf("masterContract.GetNodeHosts: %w\n", err)
	} else {
		fmt.Printf("hosts = %s\n", hosts)
	}

	contractAddr, err := masterContract.GetUserContractAddress(w.Address())
	if err != nil {
		return nil, err
	}

	contract := OpenUserContract(api, contractAddr)
	if userData, err := contract.GetData(); err != nil {
		return nil, err
	} else {
		if userData.Master.String() != masterContractAddr || userData.Owner.String() != w.Address().String() {
			return nil, errors.New("strange user contract data")
		} else {
			if contractBalance, err := contract.GetBalance(); err != nil {
				return nil, fmt.Errorf("contract.GetBalance: %w", err)
			} else {
				fmt.Printf("user contract (address = %s, balance = %s)\n", contractAddr, contractBalance)
			}
			return &UserToncli{
				api:            api,
				wallet:         w,
				masterContract: masterContract,
				contract:       contract,
			}, nil
		}
	}
}

func (c *UserToncli) GetNodeHosts() (map[string]*address.Address, error) {
	hosts, err := c.masterContract.GetNodeHosts()
	if err != nil {
		err = fmt.Errorf("masterContract.GetNodeHosts: %w", err)
	}
	return hosts, err
}

func (c *UserToncli) GetCallIds() ([]uint64, error) {
	calls, err := c.contract.GetCallIds()
	if err != nil {
		err = fmt.Errorf("GetCallIds: %w", err)
	}
	return calls, err
}

func (c *UserToncli) BuildCreateCallMessage(callId uint64) (msg, sign []byte, err error) {
	validUntil := time.Now().Unix() + 60
	msgCell := cell.BeginCell().
		MustStoreUInt(callId, 64).
		MustStoreUInt(uint64(validUntil), 32).
		EndCell()
	sign = msgCell.Sign(c.wallet.PrivateKey())
	_, msg, err = msgCell.BeginParse().RestBits()
	return
}

func (c *UserToncli) BuildEndCallMessage(callId uint64, spentMinutes uint32) (msg, sign []byte, err error) {
	validUntil := time.Now().Unix() + 60
	msgCell := cell.BeginCell().
		MustStoreUInt(callId, 64).
		MustStoreUInt(uint64(validUntil), 32).
		MustStoreUInt(uint64(spentMinutes), 32).
		EndCell()
	sign = msgCell.Sign(c.wallet.PrivateKey())
	_, msg, err = msgCell.BeginParse().RestBits()
	return
}

func (c *UserToncli) GetNodePublicKey(nodeContractAddr *address.Address) (ed25519.PublicKey, error) {
	nodeContract := OpenNodeContract(c.api, nodeContractAddr)
	if nodeContractData, err := nodeContract.GetData(); err != nil {
		return nil, fmt.Errorf("nodeContract.GetData: %w", err)
	} else {
		return nodeContractData.PublicKey, nil
	}
}

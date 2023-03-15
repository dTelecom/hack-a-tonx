package ton

import (
	"crypto/ed25519"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

const (
	opNodeWithdraw   uint64 = 0x3f6e74
	opNodeCreateCall        = 0xf3672d9
	opNodeEndCall           = 0x2c2c9c5e
)

type NodeContract struct {
	Contract
}

type NodeContractData struct {
	PublicKey ed25519.PublicKey
	NodeHost  string
	Owner     *address.Address
	Master    *address.Address
}

func OpenNodeContract(api *ton.APIClient, addr *address.Address) *NodeContract {
	return &NodeContract{
		Contract{
			api, addr,
		},
	}
}

func (c *NodeContract) SendWithdraw(via *wallet.Wallet, amount uint64) error {
	body := cell.BeginCell().
		MustStoreUInt(opNodeWithdraw, 32).
		MustStoreUInt(0, 64).
		MustStoreCoins(amount).
		EndCell()
	return c.send(via, body)
}

func (c *NodeContract) SendCreateCall(via *wallet.Wallet, userAddr *address.Address, userSign []byte, userMsg []byte) error {
	body := cell.BeginCell().
		MustStoreUInt(opNodeCreateCall, 32).
		MustStoreUInt(0, 64).
		MustStoreAddr(userAddr).
		MustStoreRef(cell.BeginCell().
			MustStoreSlice(userSign, 8*uint(len(userSign))).
			MustStoreSlice(userMsg, 8*uint(len(userMsg))).
			EndCell()).
		EndCell()
	return c.send(via, body)
}

func (c *NodeContract) SendEndCall(via *wallet.Wallet, userAddr *address.Address, userSign []byte, userMsg []byte) error {
	body := cell.BeginCell().
		MustStoreUInt(opNodeEndCall, 32).
		MustStoreUInt(0, 64).
		MustStoreAddr(userAddr).
		MustStoreRef(cell.BeginCell().
			MustStoreSlice(userSign, 8*uint(len(userSign))).
			MustStoreSlice(userMsg, 8*uint(len(userMsg))).
			EndCell()).
		EndCell()
	return c.send(via, body)
}

func (c *NodeContract) GetData() (*NodeContractData, error) {
	res, err := c.runGetMethod("get_wallet_data")
	if err != nil {
		return nil, err
	}
	return &NodeContractData{
		PublicKey: res.MustInt(0).Bytes(),
		NodeHost:  res.MustCell(1).BeginParse().MustLoadStringSnake(),
		Owner:     res.MustSlice(2).MustLoadAddr(),
		Master:    res.MustSlice(3).MustLoadAddr(),
	}, nil
}

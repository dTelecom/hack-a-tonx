package ton

import (
	"crypto/ed25519"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

const (
	opMasterWithdraw   uint64 = 0x348a7a82
	opMasterCreateUser        = 0x2b2cf99c
	opMasterCreateNode        = 0x706425c3
)

type MasterContract struct {
	Contract
}

type MasterContractData struct {
	Owner *address.Address
}

func OpenMasterContract(api *ton.APIClient, addr *address.Address) *MasterContract {
	return &MasterContract{
		Contract{
			api, addr,
		},
	}
}

func (c *MasterContract) SendWithdraw(via *wallet.Wallet, amount uint64) error {
	body := cell.BeginCell().
		MustStoreUInt(opMasterWithdraw, 32).
		MustStoreUInt(0, 64).
		MustStoreCoins(amount).
		EndCell()
	return c.send(via, body)
}

func (c *MasterContract) SendCreateUser(via *wallet.Wallet, publicKey ed25519.PublicKey) error {
	body := cell.BeginCell().
		MustStoreUInt(opMasterCreateUser, 32).
		MustStoreUInt(0, 64).
		MustStoreSlice(publicKey, 256).
		EndCell()
	return c.send(via, body)
}

func (c *MasterContract) SendCreateNode(via *wallet.Wallet, publicKey ed25519.PublicKey, nodeHost string) error {
	body := cell.BeginCell().
		MustStoreUInt(opMasterCreateNode, 32).
		MustStoreUInt(0, 64).
		MustStoreSlice(publicKey, 256).
		MustStoreRef(cell.BeginCell().MustStoreStringSnake(nodeHost).EndCell()).
		EndCell()
	return c.sendWithAmount(tlb.MustFromTON("1.1"), via, body)
}

func (c *MasterContract) GetUserContractAddress(userAddr *address.Address) (*address.Address, error) {
	param := cell.BeginCell().MustStoreAddr(userAddr).EndCell().BeginParse()
	res, err := c.runGetMethod("get_user_wallet_address", param)
	if err != nil {
		return nil, err
	}
	return res.MustSlice(0).MustLoadAddr(), nil
}

func (c *MasterContract) GetNodeContractAddress(nodeAddr *address.Address) (*address.Address, error) {
	param := cell.BeginCell().MustStoreAddr(nodeAddr).EndCell().BeginParse()
	res, err := c.runGetMethod("get_node_wallet_address", param)
	if err != nil {
		return nil, err
	}
	return res.MustSlice(0).MustLoadAddr(), nil
}

func (c *MasterContract) GetData() (*MasterContractData, error) {
	res, err := c.runGetMethod("get_dtelecom_data")
	if err != nil {
		return nil, err
	}
	return &MasterContractData{
		Owner: res.MustSlice(0).MustLoadAddr(),
	}, nil
}

func (c *MasterContract) GetNodeHosts() (map[string]*address.Address, error) {
	res, err := c.runGetMethod("get_node_hosts_list")
	if err != nil {
		return nil, err
	}

	hosts := make(map[string]*address.Address)
	for cur := res.AsTuple()[0]; cur != nil; {
		tuple := cur.([]any)
		nodeInfo := tuple[0].(*cell.Slice)
		nodeHost := nodeInfo.MustLoadRef().MustLoadStringSnake()
		nodeAddr := nodeInfo.MustLoadAddr()
		hosts[nodeHost] = nodeAddr
		cur = tuple[1]
	}
	return hosts, nil
}

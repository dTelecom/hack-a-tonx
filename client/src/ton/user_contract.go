package ton

import (
	"crypto/ed25519"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/ton"
	"math/big"
)

type UserContract struct {
	Contract
}

type UserContractData struct {
	PublicKey ed25519.PublicKey
	Owner     *address.Address
	Master    *address.Address
}

func OpenUserContract(api *ton.APIClient, addr *address.Address) *UserContract {
	return &UserContract{
		Contract{
			api, addr,
		},
	}
}

func (c *UserContract) GetData() (*UserContractData, error) {
	res, err := c.runGetMethod("get_wallet_data")
	if err != nil {
		return nil, err
	}
	return &UserContractData{
		PublicKey: res.MustInt(0).Bytes(),
		Owner:     res.MustSlice(1).MustLoadAddr(),
		Master:    res.MustSlice(2).MustLoadAddr(),
	}, nil
}

func (c *UserContract) GetCallIds() ([]uint64, error) {
	res, err := c.runGetMethod("get_call_ids_list")
	if err != nil {
		return nil, err
	}

	var calls []uint64
	for cur := res.AsTuple()[0]; cur != nil; {
		tuple := cur.([]any)
		calls = append(calls, tuple[0].(*big.Int).Uint64())
		cur = tuple[1]
	}
	return calls, nil
}

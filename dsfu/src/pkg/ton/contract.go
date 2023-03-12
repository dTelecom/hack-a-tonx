package ton

import (
	"context"
	"fmt"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

type Contract struct {
	api  *ton.APIClient
	addr *address.Address
}

func (c *Contract) send(via *wallet.Wallet, body *cell.Cell) error {
	return c.sendWithAmount(tlb.MustFromTON("0.1"), via, body)
}

func (c *Contract) sendWithAmount(amount tlb.Coins, via *wallet.Wallet, body *cell.Cell) error {
	return via.Send(context.Background(), &wallet.Message{
		Mode: 1, // pay fees separately (from balance, not from amount)
		InternalMessage: &tlb.InternalMessage{
			Bounce:  true,
			DstAddr: c.addr,
			Amount:  amount,
			Body:    body,
		},
	}, true)
}

func (c *Contract) runGetMethod(method string, params ...any) (*ton.ExecutionResult, error) {
	block, err := c.api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		return nil, fmt.Errorf("api.CurrentMasterchainInfo: %w", err)
	}
	executionResult, err := c.api.RunGetMethod(context.Background(), block, c.addr, method, params...)
	if err != nil {
		return nil, fmt.Errorf("api.RunGetMethod: %w", err)
	}
	return executionResult, nil
}

func (c *Contract) GetBalance() (tlb.Coins, error) {
	block, err := c.api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		return tlb.Coins{}, fmt.Errorf("api.CurrentMasterchainInfo: %w", err)
	}

	acc, err := c.api.GetAccount(context.Background(), block, c.addr)
	if err != nil {
		return tlb.Coins{}, fmt.Errorf("failed to get contract state: %w", err)
	}

	if !acc.IsActive {
		return tlb.Coins{}, nil
	}

	return acc.State.Balance, nil
}

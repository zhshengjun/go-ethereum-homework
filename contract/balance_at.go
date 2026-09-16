package contract

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func BalanceAt(ctx context.Context, client *ethclient.Client, addr string) *big.Int {
	balance, _ := client.BalanceAt(ctx, common.HexToAddress(addr), nil)
	return balance

}

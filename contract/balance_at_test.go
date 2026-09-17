package contract

import (
	"fmt"
	"go-ethereum-homework/internal/testutil"
	"math/big"
	"testing"
)

func TestBalanceAt(t *testing.T) {
	ctx, client, _ := testutil.IntegrationClient(t)
	address, _ := testutil.AccountFrom(t)

	fmt.Println("=================Balance Part=========================")
	fmt.Printf("address: %s\n", address)
	// 这里是对少 wei
	balance := BalanceAt(ctx, client, address)
	fmt.Println("balance(wei)", balance.String())

	// 这里需要移除 decimals，也就是除以 10^18
	ethValue := new(big.Rat).SetFrac(balance, big.NewInt(1e18))
	fmt.Println("balance(ETH)", ethValue.FloatString(3))
}

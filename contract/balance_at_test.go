package contract

import (
	"fmt"
	"go-ethereum-homework/internal/testutil"
	"math/big"
	"testing"
)

func TestBalanceAt(t *testing.T) {
	ctx, client, _ := testutil.IntegrationClient(t)

	address := "0x44b29771acb6144dDbDa47A075ed64D68Fa53c1D"

	fmt.Println("=================Balance Part=========================")
	fmt.Printf("address: %s\n", address)

	balance := BalanceAt(ctx, client, address)
	fmt.Println("balance(wei)", balance.String())
	ethValue := new(big.Rat).SetFrac(balance, big.NewInt(1e18))
	fmt.Println("balance(ETH)", ethValue.FloatString(3))

}

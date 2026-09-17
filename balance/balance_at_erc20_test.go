package balance

import (
	"go-ethereum-homework/internal/testutil"
	"math/big"
	"testing"
)

func TestBalanceAtERC20(t *testing.T) {
	ctx, client, endpoint := testutil.IntegrationClient(t)
	address, _ := testutil.AccountTo(t)

	symbol, balance, err := BalanceAtERC20(ctx, client, address, testutil.ERC20Token(t))
	if err != nil {
		t.Fatal(testutil.RedactRPCError(err, endpoint))
	}
	t.Logf("ERC20 symbol (name): %s", symbol)
	t.Logf("ERC20 balance (wei): %s", balance.String())

	// 这里需要移除 decimals，也就是除以 10^18
	ethValue := new(big.Rat).SetFrac(balance, big.NewInt(1e18))
	t.Logf("ERC20 balance: %s", ethValue.FloatString(3))
}

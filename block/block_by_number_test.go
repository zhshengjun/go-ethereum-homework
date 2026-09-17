package block

import (
	"go-ethereum-homework/internal/testutil"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

func TestLatestIntegration(t *testing.T) {
	ctx, client, endpoint := testutil.IntegrationClient(t)

	chainID, err := client.ChainID(ctx)
	if err != nil {
		t.Fatal(strings.ReplaceAll(err.Error(), endpoint, "<rpc-url>"))
	}

	// 这里 number不传值，如果是 nil 就表示 latest
	latest, err := ByNumber(ctx, client, big.NewInt(int64(rpc.LatestBlockNumber)))

	if err != nil {
		t.Fatal(strings.ReplaceAll(err.Error(), endpoint, "<rpc-url>"))
	}
	if latest.Number() == nil {
		t.Fatal("latest block number is nil")
	}

	header := latest.Header()
	t.Logf("RPC URL: %s", endpoint)
	t.Logf("Chain ID: %s", chainID)
	t.Logf("Block Number: %s", latest.Number())
	t.Logf("Block Hash: %s", latest.Hash())
	t.Logf("Header Hash: %s", header.Hash())
	t.Logf("Parent Hash: %s", header.ParentHash)
	t.Logf("Withdrawals: %d", len(latest.Withdrawals()))
	t.Logf("Gas BaseFee: %d", latest.BaseFee())
	t.Logf("Gas Limit: %d", latest.GasLimit())
	t.Logf("Gas Used: %d", latest.GasUsed())
	t.Logf("State Root: %s", header.Root)
	t.Logf("Block Time: %s", time.Unix(int64(header.Time), 0).Format(time.DateTime))
}

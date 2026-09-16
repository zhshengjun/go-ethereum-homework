package block

import (
	"go-ethereum-homework/internal/testutil"
	"strings"
	"testing"
	"time"
)

func TestLatestIntegration(t *testing.T) {
	ctx, client, endpoint := testutil.IntegrationClient(t)

	chainID, err := client.ChainID(ctx)
	if err != nil {
		t.Fatal(strings.ReplaceAll(err.Error(), endpoint, "<rpc-url>"))
	}

	latest, err := ByNumber(ctx, client, nil)
	if err != nil {
		t.Fatal(strings.ReplaceAll(err.Error(), endpoint, "<rpc-url>"))
	}
	if latest.Number() == nil {
		t.Fatal("latest block number is nil")
	}

	header := latest.Header()
	t.Log("RPC URL: configured via config.yaml rpc.url")
	t.Logf("Chain ID: %s", chainID)
	t.Logf("Block Number: %s", latest.Number())
	t.Logf("Block Hash: %s", latest.Hash())
	t.Logf("Header Hash: %s", header.Hash())
	t.Logf("Parent Hash: %s", header.ParentHash)
	t.Logf("Withdrawals: %d", len(latest.Withdrawals()))
	t.Logf("Gas Limit: %d", latest.GasLimit())
	t.Logf("Gas Used: %d", latest.GasUsed())
	t.Logf("State Root: %s", header.Root)
	if header.SlotNumber == nil {
		t.Log("Slot Number: unavailable")
	} else {
		t.Logf("Slot Number: %d", *header.SlotNumber)
	}
	t.Logf("Block Time: %s", time.Unix(int64(header.Time), 0).Format(time.RFC3339))
}

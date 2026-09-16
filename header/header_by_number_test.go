package header

import (
	"go-ethereum-homework/internal/testutil"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
)

func TestHeaderByNumberIntegration(t *testing.T) {
	ctx, client, endpoint := testutil.IntegrationClient(t)

	latest, err := ByNumber(ctx, client, nil)
	if err != nil {
		t.Fatal(testutil.RedactRPCError(err, endpoint))
	}
	safe, err := ByNumber(ctx, client, big.NewInt(int64(rpc.SafeBlockNumber)))
	if err != nil {
		t.Fatal(testutil.RedactRPCError(err, endpoint))
	}
	finalized, err := ByNumber(ctx, client, big.NewInt(int64(rpc.FinalizedBlockNumber)))
	if err != nil {
		t.Fatal(testutil.RedactRPCError(err, endpoint))
	}

	if finalized.Number.Cmp(safe.Number) > 0 || safe.Number.Cmp(latest.Number) > 0 {
		t.Fatalf("unexpected order: finalized=%s safe=%s latest=%s", finalized.Number, safe.Number, latest.Number)
	}

	t.Logf("Latest: %s %s", latest.Number, latest.Hash())
	t.Logf("Safe: %s %s", safe.Number, safe.Hash())
	t.Logf("Finalized: %s %s", finalized.Number, finalized.Hash())
}

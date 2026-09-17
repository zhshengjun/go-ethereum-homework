package header

import (
	"go-ethereum-homework/internal/testutil"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
)

func TestHeaderByHashIntegration(t *testing.T) {
	ctx, client, endpoint := testutil.IntegrationClient(t)

	saferBlock, err := ByNumber(ctx, client, big.NewInt(int64(rpc.SafeBlockNumber)))
	if err != nil {
		t.Fatal(err, endpoint)
	}

	got, err := ByHash(ctx, client, saferBlock.Hash())
	if err != nil {
		t.Fatal(err, endpoint)
	}
	if got.Hash() != saferBlock.Hash() {
		t.Fatalf("hash mismatch: got=%s want=%s", got.Hash(), saferBlock.Hash())
	}

	t.Logf("Block Number: %s", got.Number)
	t.Logf("Block Hash: %s", got.Hash())
	t.Logf("Parent Hash: %s", got.ParentHash)
}

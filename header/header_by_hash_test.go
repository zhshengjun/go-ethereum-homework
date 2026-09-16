package header

import (
	"go-ethereum-homework/internal/testutil"
	"testing"
)

func TestHeaderByHashIntegration(t *testing.T) {
	ctx, client, endpoint := testutil.IntegrationClient(t)

	latest, err := ByNumber(ctx, client, nil)
	if err != nil {
		t.Fatal(err, endpoint)
	}

	got, err := ByHash(ctx, client, latest.Hash())
	if err != nil {
		t.Fatal(err, endpoint)
	}
	if got.Hash() != latest.Hash() {
		t.Fatalf("hash mismatch: got=%s want=%s", got.Hash(), latest.Hash())
	}

	t.Logf("Block Number: %s", got.Number)
	t.Logf("Block Hash: %s", got.Hash())
	t.Logf("Parent Hash: %s", got.ParentHash)
}

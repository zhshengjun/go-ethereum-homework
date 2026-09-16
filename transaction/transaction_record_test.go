package transaction

import (
	"go-ethereum-homework/internal/testutil"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestRecord(t *testing.T) {
	ctx, client, _ := testutil.IntegrationWSSClient(t)
	txHash := testutil.TransactionHash(t)
	if common.HexToHash(txHash) == (common.Hash{}) {
		t.Skip("set transaction.hash in config.yaml")
	}
	Record(ctx, client, txHash)
}

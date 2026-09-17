package transaction

import (
	"context"
	"testing"
	"time"

	"go-ethereum-homework/internal/testutil"
)

func TestSendTransactionBind(t *testing.T) {

	_, client, _ := testutil.IntegrationWSSClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	_, privateKey := testutil.AccountFrom(t)
	address, _ := testutil.AccountTo(t)
	receipt, err := TransactionBind(
		ctx,
		client,
		privateKey,
		address,
		0.001,
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("TxHash: %s", receipt.TxHash.Hex())
	t.Logf("Status: %d", receipt.Status)
	t.Logf("BlockNumber: %s", receipt.BlockNumber)
	t.Logf("BlockHash: %s", receipt.BlockHash.Hex())
	t.Logf("TransactionIndex: %d", receipt.TransactionIndex)
	t.Logf("GasUsed: %d", receipt.GasUsed)
	t.Logf("EffectiveGasPrice: %s wei", receipt.EffectiveGasPrice)
	t.Logf("Logs: %d", len(receipt.Logs))
}

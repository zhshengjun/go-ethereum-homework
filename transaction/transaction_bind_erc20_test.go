package transaction

import (
	"go-ethereum-homework/internal/testutil"
	"testing"
)

func TestSendERC20(t *testing.T) {

	ctx, client, _ := testutil.IntegrationWSSClient(t)
	_, privateKey := testutil.AccountFrom(t)
	address, _ := testutil.AccountTo(t)
	token := testutil.ERC20Token(t)
	receipt := SendERC20(
		ctx,
		client,
		privateKey,
		address,
		token,
		5*1e18,
	)

	t.Logf("TxHash: %s", receipt.TxHash.Hex())
	t.Logf("Status: %d", receipt.Status)
	t.Logf("BlockNumber: %s", receipt.BlockNumber)
	t.Logf("BlockHash: %s", receipt.BlockHash.Hex())
	t.Logf("TransactionIndex: %d", receipt.TransactionIndex)
	t.Logf("GasUsed: %d", receipt.GasUsed)
	t.Logf("EffectiveGasPrice: %s wei", receipt.EffectiveGasPrice)
	t.Logf("Logs: %d", len(receipt.Logs))
}

package transaction

import (
	"go-ethereum-homework/internal/testutil"
	"testing"
)

func TestSendTransaction(t *testing.T) {

	ctx, client, _ := testutil.IntegrationWSSClient(t)
	_, privateKey := testutil.AccountFrom(t)
	address, _ := testutil.AccountTo(t)
	SendTransaction(
		ctx,
		client,
		privateKey,
		address,
		0.0001,
		nil,
	)
}

package transaction

import (
	"go-ethereum-homework/internal/testutil"
	"testing"
)

func TestSendTransaction(t *testing.T) {

	ctx, client, _ := testutil.IntegrationWSSClient(t)
	address, privateKey := testutil.Accoun(t)
	SendTransaction(
		ctx,
		client,
		privateKey,
		address,
		0.001,
		nil,
	)
}

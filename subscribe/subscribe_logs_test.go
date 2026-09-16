package subscribe

import (
	"context"
	"errors"
	"go-ethereum-homework/internal/testutil"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestSubscribeLogsIntegration(t *testing.T) {
	ctx, client, endpoint := testutil.IntegrationWSSClient(t)

	usdc := common.HexToAddress("0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238")
	code, err := client.CodeAt(ctx, usdc, nil)
	if err != nil {
		t.Fatal(strings.ReplaceAll(err.Error(), endpoint, "<rpc-wss>"))
	}
	if len(code) == 0 {
		t.Fatalf("contract %s is not deployed on the configured network", usdc)
	}

	err = SubscribeLogs(ctx, client, usdc, 0)
	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		t.Fatal(strings.ReplaceAll(err.Error(), endpoint, "<rpc-wss>"))
	}
}

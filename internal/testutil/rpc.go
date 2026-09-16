package testutil

import (
	"context"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

func RPCURL(t testing.TB) string {
	return rpcValue(t, "rpc.url")
}

func RPCWSS(t testing.TB) string {
	return rpcValue(t, "rpc.wss")
}

func rpcValue(t testing.TB, key string) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate config.yaml")
	}

	config := viper.New()
	config.SetConfigFile(filepath.Join(filepath.Dir(file), "..", "..", "config.yaml"))
	if err := config.ReadInConfig(); err != nil {
		t.Fatal(err)
	}

	endpoint := config.GetString(key)
	if endpoint == "" {
		t.Fatalf("config.yaml %s is empty", key)
	}
	return endpoint
}

func IntegrationClient(t *testing.T) (context.Context, *ethclient.Client, string) {
	t.Helper()

	endpoint := RPCURL(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)

	client, err := ethclient.DialContext(ctx, endpoint)
	if err != nil {
		t.Fatal(RedactRPCError(err, endpoint))
	}
	t.Cleanup(client.Close)

	return ctx, client, endpoint
}

func IntegrationWSSClient(t *testing.T) (context.Context, *ethclient.Client, string) {
	t.Helper()

	endpoint := RPCWSS(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	t.Cleanup(cancel)

	client, err := ethclient.DialContext(ctx, endpoint)
	if err != nil {
		t.Fatal(RedactRPCError(err, endpoint))
	}
	t.Cleanup(client.Close)

	return ctx, client, endpoint
}

func RedactRPCError(err error, endpoint string) string {
	return strings.ReplaceAll(err.Error(), endpoint, "<rpc-url>")
}

func TransactionHash(t testing.TB) string {
	return rpcValue(t, "transaction.hash")
}

func Accoun(t testing.TB) (string, string) {
	address := rpcValue(t, "account.to.address")
	private := rpcValue(t, "account.from.private")

	return address, private
}

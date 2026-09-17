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
	return configValue(t, "rpc.url")
}

func RPCWSS(t testing.TB) string {
	return configValue(t, "rpc.wss")
}

// configValue 获取配置值
func configValue(t testing.TB, key string) string {
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
	return integrationClient(t, RPCURL, 10*time.Second)
}

func IntegrationWSSClient(t *testing.T) (context.Context, *ethclient.Client, string) {
	t.Helper()
	return integrationClient(t, RPCWSS, 15*time.Second)
}

func integrationClient(
	t *testing.T,
	endpointFunc func(testing.TB) string,
	timeout time.Duration,
) (context.Context, *ethclient.Client, string) {
	t.Helper()

	endpoint := endpointFunc(t)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	// 这里现注册清理逻辑，测试结束后会调用函数
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

// TransactionHashHex 获取交易hash
func TransactionHashHex(t testing.TB) string {
	return configValue(t, "transaction.hash")
}

// AccountFrom 获取账号
func AccountFrom(t testing.TB) (string, string) {
	address := configValue(t, "account.from.address")
	private := configValue(t, "account.from.private")

	return address, private
}

func AccountTo(t testing.TB) (string, string) {
	address := configValue(t, "account.to.address")
	private := configValue(t, "account.to.private")

	return address, private
}

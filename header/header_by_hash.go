package header

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func ByHash(ctx context.Context, client *ethclient.Client, hash common.Hash) (*types.Header, error) {
	return client.HeaderByHash(ctx, hash)
}

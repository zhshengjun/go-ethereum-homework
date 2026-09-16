package header

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func ByNumber(ctx context.Context, client *ethclient.Client, number *big.Int) (*types.Header, error) {
	return client.HeaderByNumber(ctx, number)
}

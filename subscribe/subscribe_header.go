package subscribe

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Subscribe(
	ctx context.Context,
	client *ethclient.Client,
	maxHeaders int,
) error {
	headers := make(chan *types.Header)

	sub, err := client.SubscribeNewHead(ctx, headers)
	if err != nil {
		return fmt.Errorf("subscribe new head: %w", err)
	}
	defer sub.Unsubscribe()

	received := 0

	for {
		select {
		case err, ok := <-sub.Err():
			if !ok {
				return nil
			}
			if err != nil {
				return fmt.Errorf("subscription error: %w", err)
			}

		case header, ok := <-headers:
			if !ok {
				return nil
			}

			if header == nil {
				continue
			}

			fmt.Printf("receiver Header %s\n", header.Number.String())
			fmt.Print(formatHeader(header))

			received++

			// maxHeaders <= 0 表示持续监听
			if maxHeaders > 0 && received >= maxHeaders {
				return nil
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func formatHeader(header *types.Header) string {
	return fmt.Sprintf(
		"Header Hash: %s\nHeader Number: %d\nHeader ParentHash: %s\nHeader GasUsed: %d\nHeader Time: %s\n",
		header.Hash().Hex(),
		header.Number,
		header.ParentHash.Hex(),
		header.GasUsed,
		time.Unix(int64(header.Time), 0).Format(time.DateTime),
	)
}

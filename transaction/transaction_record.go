package transaction

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Record(ctx context.Context, client *ethclient.Client, txHash string) {

	tx, _, err := client.TransactionByHash(ctx, common.HexToHash(txHash))
	if err != nil {
		fmt.Println("get transaction error:", err)
		return
	}
	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		fmt.Println("recover sender error:", err)
		return
	}

	fmt.Println("=================Transaction Part=========================")
	fmt.Printf("Transaction From: %s\n", from)
	if tx.To() == nil {
		fmt.Println("Transaction To: contract creation")
	} else {
		fmt.Printf("Transaction To: %s\n", tx.To())
	}
	fmt.Printf("Transaction Nonce: %d\n", tx.Nonce())
	fmt.Printf("Transaction Time: %s\n", tx.Time().Format(time.DateTime))
	fmt.Printf("Transaction Amount（wei）: %s\n", tx.Value())
	fmt.Printf("Transaction Input: %s\n", tx.Data())
	fmt.Printf("Gas Fee Cap: %s\n", tx.GasFeeCap())
	fmt.Printf("Gas Price: %s\n", tx.GasPrice())
	fmt.Printf("Gas Tip Cap: %s\n", tx.GasTipCap())
	fmt.Printf("Cost = Amount + GasFeeCap * GasLimit: %s\n", tx.Cost())

	receipt, err := client.TransactionReceipt(ctx, common.HexToHash(txHash))
	if err != nil {
		fmt.Println("get receipt error:", err)
		return
	}
	fmt.Println("=================Receipt Part=========================")
	fmt.Printf("Receipt Status: %d\n", receipt.Status)
	fmt.Printf("Receipt BlockNumber: %s\n", receipt.BlockNumber)
	fmt.Printf("Receipt BlockHash: %s\n", receipt.BlockHash)
	fmt.Printf("Receipt ContractAddress: %s\n", receipt.ContractAddress)
	fmt.Printf("Receipt GasUsed: %d\n", receipt.GasUsed)
	fmt.Printf("Receipt EffectiveGasPrice: %s\n", receipt.EffectiveGasPrice)

	gasPrice := new(big.Int).Mul(new(big.Int).SetUint64(receipt.GasUsed), receipt.EffectiveGasPrice)
	fmt.Printf("GasUsed * EffectiveGasPrice(wei): %s\n", gasPrice)
	fmt.Printf("GasUsed * EffectiveGasPrice(ETH): %s\n", new(big.Rat).SetFrac(gasPrice, big.NewInt(1e18)).FloatString(18))

	if receipt.BlobGasUsed > 0 {
		fmt.Printf("Receipt BlobGasUsed: %d\n", receipt.BlobGasUsed)
		fmt.Printf("Receipt BlobGasPrice: %s\n", receipt.BlobGasPrice)
		fmt.Printf("Blob GasUsed: %s\n", new(big.Int).Mul(new(big.Int).SetUint64(receipt.BlobGasUsed), receipt.BlobGasPrice))
	}
}

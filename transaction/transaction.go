package transaction

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// SendTransaction 测试再go中使用 client 发送交易
func SendTransaction(ctx context.Context, client *ethclient.Client,
	private string, to string, amount float64, input []byte) {
	// 目标地址
	toAddress := common.HexToAddress(to)
	// 私钥
	privateKey, err := crypto.HexToECDSA(private)
	if err != nil {
		log.Fatal(err)
	}
	from := crypto.PubkeyToAddress(privateKey.PublicKey)

	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatal(err)
	}
	nonce, err := client.PendingNonceAt(ctx, from)
	if err != nil {
		log.Fatal(err)
	}
	gasTipCap, err := client.SuggestGasTipCap(ctx)
	if err != nil {
		log.Fatal(err)
	}
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	baseFee := header.BaseFee
	if baseFee == nil {
		baseFee, err = client.SuggestGasPrice(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}
	gasFeeCap := new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), gasTipCap)

	amountWei, _ := new(big.Float).Mul(big.NewFloat(amount), big.NewFloat(1e18)).Int(new(big.Int))

	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
		From:  from,
		To:    &toAddress,
		Value: amountWei,
		Data:  input,
	})
	if err != nil {
		log.Fatal(err)
	}

	totalCost := new(big.Int).Add(amountWei, new(big.Int).Mul(gasFeeCap, new(big.Int).SetUint64(gasLimit)))
	balance, err := client.BalanceAt(ctx, from, nil)
	if err != nil {
		log.Fatal(err)
	}
	if balance.Cmp(totalCost) < 0 {
		fmt.Println("account balance is less than total cost")
		return
	}

	txData := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gasLimit,
		To:        &toAddress,
		Value:     amountWei,
		Data:      input,
	}

	tx := types.NewTx(txData)
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
	if err != nil {
		log.Fatal("SignTx:", err)
	}
	if err := client.SendTransaction(ctx, signedTx); err != nil {
		log.Fatal("Send tx:", err)
	}

	fmt.Printf("From: %s\n", from)
	fmt.Printf("To: %s\n", toAddress)
	fmt.Printf("Value: %.6f ETH (%s wei)\n", amount, amountWei)
	fmt.Printf("Gas Limit: %d\n", gasLimit)
	fmt.Printf("Gas Tip Cap: %s\n", gasTipCap)
	fmt.Printf("Gas Fee Cap: %s\n", gasFeeCap)
	fmt.Printf("Total Cost: %s\n", totalCost)
	fmt.Printf("Nonce: %d\n", nonce)
	fmt.Printf("Tx Hash: %s\n", signedTx.Hash())
}

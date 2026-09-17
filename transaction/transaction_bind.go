package transaction

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TransactionBind 测试再go中使用 client 发送交易
func TransactionBind(ctx context.Context, client *ethclient.Client,
	private string, to string, amount float64) (*types.Receipt, error) {
	// 发送方的私钥,获取发送的地址
	privateKey, err := crypto.HexToECDSA(private)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get chain ID: %w", err)
	}

	transactOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("create transactor: %w", err)
	}

	amountWei, _ := new(big.Float).Mul(big.NewFloat(amount), big.NewFloat(1e18)).Int(new(big.Int))

	toAddress := common.HexToAddress(to)
	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
		From:  transactOpts.From,
		To:    &toAddress,
		Value: amountWei,
	})
	if err != nil {
		return nil, fmt.Errorf("estimate gas: %w", err)
	}

	transactOpts.Context = ctx
	transactOpts.Value = amountWei
	transactOpts.GasLimit = gasLimit

	receiver := bind.NewBoundContract(
		toAddress,
		abi.ABI{},
		nil,
		client,
		nil,
	)

	//balance, _ :=client.PendingBalanceAt(ctx, crypto.PubkeyToAddress(privateKey.PublicKey))

	tx, err := receiver.RawTransact(transactOpts, nil)
	// 如果有abi，相当于ERC20token了
	// tx, err := receiver.Transact(transactOpts, "transfer", to, amountWei)
	if err != nil {
		return nil, fmt.Errorf("send transaction: %w", err)
	}
	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return nil, fmt.Errorf("wait mined: %w", err)
	}

	return receipt, nil
}

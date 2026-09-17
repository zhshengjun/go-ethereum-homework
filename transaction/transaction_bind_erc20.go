package transaction

import (
	"context"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const erc20ABI = `[
  {
    "type": "function",
    "name": "transfer",
    "stateMutability": "nonpayable",
    "inputs": [
      {"name": "to", "type": "address"},
      {"name": "value", "type": "uint256"}
    ],
    "outputs": [
      {"name": "", "type": "bool"}
    ]
  }
]`

// SendTransaction 测试再go中使用 client 发送交易
func SendERC20(ctx context.Context, client *ethclient.Client,
	private string, to string, tokenAddr string, amount uint64) *types.Receipt {
	privateKey, _ := crypto.HexToECDSA(private)

	chainID, _ := client.ChainID(ctx)
	transactOpts, _ := bind.NewKeyedTransactorWithChainID(privateKey, chainID)

	// 获取合约
	parsedABI, _ := abi.JSON(strings.NewReader(erc20ABI))
	token := bind.NewBoundContract(
		common.HexToAddress(tokenAddr),
		parsedABI,
		client,
		client,
		client,
	)

	// 发起转账
	tx, _ := token.Transact(
		transactOpts,
		"transfer",
		common.HexToAddress(to),
		new(big.Int).SetUint64(amount),
	)

	receipt, _ := bind.WaitMined(ctx, client, tx)

	return receipt

}

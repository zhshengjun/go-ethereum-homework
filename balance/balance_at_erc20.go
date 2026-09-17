package balance

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const erc20BalanceABI = `
[
  {
    "type": "function",
    "name": "balanceOf",
    "stateMutability": "view",
    "inputs": [
      {
        "name": "account",
        "type": "address"
      }
    ],
    "outputs": [
      {
        "name": "",
        "type": "uint256"
      }
    ]
  },
  {
    "type": "function",
    "name": "symbol",
    "inputs": [],
    "outputs": [
      {
        "name": "",
        "type": "string",
        "internalType": "string"
      }
    ],
    "stateMutability": "view"
  }
]
`

func BalanceAtERC20(ctx context.Context, client *ethclient.Client, addr string, tokenAddr string) (string, *big.Int, error) {
	if !common.IsHexAddress(addr) || !common.IsHexAddress(tokenAddr) {
		return "", nil, fmt.Errorf("invalid account or token address")
	}
	// 解析abi
	parsedABI, err := abi.JSON(strings.NewReader(erc20BalanceABI))
	if err != nil {
		return "", nil, err
	}

	// 构造合约
	token := bind.NewBoundContract(common.HexToAddress(tokenAddr), parsedABI, client, nil, nil)
	symbolValues, err := getValue(ctx, token, "symbol")
	if err != nil {
		return "", nil, err
	}
	symbol, ok := symbolValues[0].(string)
	if !ok {
		return "", nil, fmt.Errorf("symbol returned %T", symbolValues[0])
	}

	balanceValues, err := getValue(ctx, token, "balanceOf", common.HexToAddress(addr))
	if err != nil {
		return "", nil, err
	}
	balance, ok := balanceValues[0].(*big.Int)
	if !ok {
		return "", nil, fmt.Errorf("balanceOf returned %T", balanceValues[0])
	}
	return symbol, balance, nil

}

func getValue(ctx context.Context, token *bind.BoundContract, method string, params ...any) ([]any, error) {

	var values []any
	if err := token.Call(&bind.CallOpts{Context: ctx}, &values, method, params...); err != nil {
		return nil, err
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("%s returned %d values", method, len(values))
	}
	return values, nil
}

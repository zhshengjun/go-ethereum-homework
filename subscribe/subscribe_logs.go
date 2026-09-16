package subscribe

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const erc20ABI_Json = `
[
	{
		"anonymous": false,
		"inputs": [
			{"indexed": true, "name": "from", "type":"address"},
			{"indexed": true, "name": "to", "type":"address"},
			{"indexed": false, "name": "value", "type":"uint256"}
		],
			"name": "Transfer",
			"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{"indexed": true, "name": "owner", "type":"address"},
			{"indexed": true, "name": "spender", "type":"address"},
			{"indexed": false, "name": "value", "type":"uint256"}
		],
		"name": "Approval",
		"type": "event"
	}
]
`

func SubscribeLogs(
	ctx context.Context,
	client *ethclient.Client,
	contractAddress common.Address,
	maxLogs int,
) error {

	parsedABI, err := abi.JSON(strings.NewReader(erc20ABI_Json))
	if err != nil {
		return err
	}

	transferEvent, ok := parsedABI.Events["Transfer"]
	if !ok {
		return fmt.Errorf("Transfer event not found")
	}

	fmt.Printf(
		"subscribe to Transfer logs of contract: %s\n",
		contractAddress.Hex(),
	)

	logsCh := make(chan types.Log)

	filter := map[string]interface{}{
		"address": contractAddress,
		"topics": [][]common.Hash{
			{transferEvent.ID},
		},
	}

	// 直接使用底层 RPC，不让 ethclient 自动加
	// fromBlock: 0x0
	// toBlock: latest
	sub, err := client.Client().EthSubscribe(
		ctx,
		logsCh,
		"logs",
		filter,
	)
	if err != nil {
		return fmt.Errorf("EthSubscribe failed: %w", err)
	}

	defer sub.Unsubscribe()

	received := 0

	for {
		select {

		case <-ctx.Done():
			return ctx.Err()

		case err, ok := <-sub.Err():
			if !ok {
				return nil
			}
			return fmt.Errorf("subscription error: %w", err)

		case log, ok := <-logsCh:
			if !ok {
				return nil
			}

			printfLog(log, parsedABI)

			received++

			if maxLogs > 0 && received >= maxLogs {
				return nil
			}
		}
	}
}

func printfLog(subscribe types.Log, json abi.ABI) {
	// 如果没有，则跳过
	if len(subscribe.Topics) == 0 {
		return
	}
	// 识别类型
	// Topics[x] 的值是事件签名的 keccak256的结果
	// Transfer(address,address,uint256)
	eventName, eventSig, ok := eventByTopic(subscribe.Topics[0], json)
	if !ok {
		fmt.Println("No event topic")
		return
	}

	// 解析参数
	fmt.Printf("event: %s(%s)\n", eventName, time.Now().Format(time.DateTime))
	fmt.Printf("BlockNumber: %d\n", subscribe.BlockNumber)
	fmt.Printf("Tx Hash : %s\n", subscribe.TxHash.Hex())
	fmt.Printf("Tx Index : %d\n", subscribe.TxIndex)
	fmt.Printf("Log Index : %d\n", subscribe.Index) // 第 X 笔交易

	values, _ := json.Unpack(eventName, subscribe.Data)

	nonIndexedInputs := make([]abi.Argument, 0)
	for _, input := range eventSig.Inputs {
		if !input.Indexed {
			nonIndexedInputs = append(nonIndexedInputs, input)
		}
	}

	printEvent(subscribe, eventSig, values)
}

func eventByTopic(topic common.Hash, parsedABI abi.ABI) (string, abi.Event, bool) {
	for name, event := range parsedABI.Events {
		if crypto.Keccak256Hash([]byte(event.Sig)) == topic {
			return name, event, true
		}
	}
	return "", abi.Event{}, false
}

// Topics[0] = Transfer/Approved 事件签名
// Topics[1] = from /owner
// Topics[2] = to spender
// Data      = value

// eventSig.Inputs[0] -> subscribe.Topics[1] // from / owner
// eventSig.Inputs[1] -> subscribe.Topics[2] // to / spender
// eventSig.Inputs[2] -> subscribe.Data      // value
func printEvent(subscribe types.Log, eventSig abi.Event, values []any) {
	topicIndex := 1 // Topics[0] 是事件签名
	dataIndex := 0  // data 的索引，正常情况就是1个

	if eventSig.Anonymous {
		topicIndex = 0
	}

	for i, input := range eventSig.Inputs {
		if input.Indexed {
			if topicIndex >= len(subscribe.Topics) {
				continue
			}

			topic := subscribe.Topics[topicIndex]

			switch input.Type.T {
			case abi.AddressTy:
				fmt.Printf(" [%d],%s(%s): %s\n",
					i, input.Name, input.Type,
					common.BytesToAddress(topic.Bytes()).Hex(),
				)
			case abi.UintTy, abi.IntTy:
				fmt.Printf(" [%d],%s(%s): %s\n",
					i, input.Name, input.Type,
					new(big.Int).SetBytes(topic.Bytes()).String(),
				)
			default:
				fmt.Printf(" [%d],%s(%s): %s\n",
					i, input.Name, input.Type, topic.Hex(),
				)
			}

			topicIndex++
			continue
		}

		if dataIndex >= len(values) {
			continue
		}

		switch value := values[dataIndex].(type) {
		case *big.Int:
			fmt.Printf(" [%d],%s(%s): %s\n",
				i, input.Name, input.Type, value.String(),
			)
		default:
			fmt.Printf(" [%d],%s(%s): %v\n",
				i, input.Name, input.Type, value,
			)
		}

		dataIndex++
	}
}

// Topics[0] = Transfer/Approved 事件签名
// Topics[1] = from /owner
// Topics[2] = to spender
// Data      = value
// eventSig.Inputs[0] -> subscribe.Topics[1] // from / owner
// eventSig.Inputs[1] -> subscribe.Topics[2] // to / spender
// eventSig.Inputs[2] -> subscribe.Data      // value
func printIndexed(subscribe types.Log, eventSig abi.Event, indexed *bool) {
	for i := range 3 {
		input := eventSig.Inputs[i]
		if indexed != nil && input.Indexed != *indexed {
			continue
		}
		topic := subscribe.Topics[i+1]

		switch input.Type.T {
		case abi.AddressTy:
			addr := common.BytesToAddress(topic.Bytes())
			fmt.Printf(" [%d],%s(%s):%s\n", i, input.Name, input.Type, addr.Hex())
		case abi.IntTy, abi.UintTy:
			value := new(big.Int).SetBytes(topic.Bytes())
			fmt.Printf(" [%d],%s(%s):%s\n", i, input.Name, input.Type, value)
		case abi.BoolTy:
			fmt.Printf("  %t\n", topic[31] != 0)
		default:
			fmt.Printf(" Unknown type: %s\n", topic.Hex())
		}
	}
}

// Topics[0] = Transfer/Approved 事件签名
// Topics[1] = from /owner
// Topics[2] = to spender
// Data      = value
// 这里就是将 交易记录中的 data 中的数据输出
func printNonIndexed(subscribe types.Log, eventSig abi.Event, values []any) {
	if len(subscribe.Data) == 0 {
		return
	}
	nonIndexedInputs := make([]abi.Argument, 0)
	for _, input := range eventSig.Inputs {
		if !input.Indexed {
			nonIndexedInputs = append(nonIndexedInputs, input)
		}
	}

	if len(nonIndexedInputs) == 0 {
		return
	}
	nonIndexedIdx := 0
	for i, input := range eventSig.Inputs {
		if !input.Indexed {
			if nonIndexedIdx < len(nonIndexedInputs) {
				value := values[nonIndexedIdx]
				switch t := value.(type) {
				case *big.Int:
					fmt.Printf(" [%d],%s(%s): %s\n", i, input.Name, input.Type, t.String())
				case common.Address:
					fmt.Printf("  %s\n", t.Hex())
				case []byte:
					fmt.Printf("  %x\n", t)
				default:
					fmt.Printf("Unknown type: %v\n", t)
				}
				nonIndexedIdx++
			}
		}
	}

}

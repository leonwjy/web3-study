package lesson

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// WebSocket 版本 - 实时订阅
func subscribeNewHeads() {
	projectID := os.Getenv("INFURA_PROJECT_ID")
	if projectID == "" {
		projectID = "64820e5ca95c49c9b5022409694643ce" // 临时使用，生产环境请用环境变量
	}

	client, err := ethclient.Dial(fmt.Sprintf("wss://sepolia.infura.io/ws/v3/%s", projectID))
	if err != nil {
		log.Fatal(err)
	}

	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("开始监听新区块...")

	// 只监听一个区块就退出（用于测试）
	select {
	case err := <-sub.Err():
		log.Fatal(err)
	case header := <-headers:
		fmt.Println("收到新区块!")
		printBlockInfo(header, client)
	case <-time.After(10 * time.Second): // 10秒超时，用于测试
		fmt.Println("10秒内没有收到新区块，使用当前最新区块进行测试")
		// 获取当前最新区块进行测试
		latestHeader, err := client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			log.Fatal("获取最新区块头失败:", err)
		}
		printBlockInfo(latestHeader, client)
	}

	sub.Unsubscribe() // 取消订阅
}

// HTTP 轮询版本 - 定期检查新区块
func pollNewHeads() {
	projectID := os.Getenv("INFURA_PROJECT_ID")
	if projectID == "" {
		projectID = "64820e5ca95c49c9b5022409694643ce" // 临时使用，生产环境请用环境变量
	}

	client, err := ethclient.Dial(fmt.Sprintf("https://sepolia.infura.io/v3/%s", projectID))
	if err != nil {
		log.Fatal(err)
	}

	var lastBlockNumber uint64 = 0

	for {
		header, err := client.HeaderByNumber(context.Background(), nil) // 获取最新区块头
		if err != nil {
			log.Printf("获取区块头失败: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		currentNumber := header.Number.Uint64()
		if currentNumber > lastBlockNumber {
			fmt.Printf("发现新区块: %d\n", currentNumber)
			printBlockInfo(header, client)
			lastBlockNumber = currentNumber
		}

		time.Sleep(12 * time.Second) // Sepolia 出块时间约12秒
	}
}

func printBlockInfo(header *types.Header, client *ethclient.Client) {
	fmt.Println("区块哈希:", header.Hash().Hex())
	fmt.Println("区块号:", header.Number.Uint64())

	// 重试机制：尝试多次获取完整区块，因为区块头可能比区块体先到达
	var block *types.Block
	var err error
	for retries := 0; retries < 3; retries++ {
		block, err = client.BlockByHash(context.Background(), header.Hash())
		if err == nil {
			break
		}
		fmt.Printf("获取区块失败 (尝试 %d/3): %v\n", retries+1, err)
		time.Sleep(2 * time.Second) // 等待2秒后重试
	}

	if err != nil {
		fmt.Printf("最终获取区块失败，显示区块头信息:\n")
		fmt.Printf("  区块号: %d\n", header.Number.Uint64())
		fmt.Printf("  难度: %d\n", header.Difficulty.Uint64())
		fmt.Printf("  Gas Limit: %d\n", header.GasLimit)
		fmt.Printf("  Gas Used: %d\n", header.GasUsed)
		fmt.Printf("  时间戳: %d\n", header.Time)
		fmt.Printf("  交易根哈希: %s\n", header.TxHash.Hex())
		fmt.Println("---")
		return
	}

	// 获取成功，打印完整信息
	fmt.Println("完整区块信息:")
	fmt.Println("  区块哈希:", block.Hash().Hex())
	fmt.Println("  区块号:", block.Number().Uint64())
	fmt.Println("  时间戳:", block.Time())
	fmt.Println("  Nonce:", block.Nonce())
	fmt.Println("  交易数量:", len(block.Transactions()))
	fmt.Println("  Gas Used:", block.GasUsed())
	fmt.Println("  Gas Limit:", block.GasLimit())
	fmt.Println("---")
}

// 测试连接和获取最新区块
func testConnection() {
	projectID := os.Getenv("INFURA_PROJECT_ID")
	if projectID == "" {
		projectID = "64820e5ca95c49c9b5022409694643ce"
	}

	client, err := ethclient.Dial(fmt.Sprintf("https://sepolia.infura.io/v3/%s", projectID))
	if err != nil {
		log.Fatal("连接失败:", err)
	}

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal("获取区块头失败:", err)
	}

	fmt.Println("连接成功!")
	fmt.Println("最新区块号:", header.Number.Uint64())
	fmt.Println("最新区块哈希:", header.Hash().Hex())
	fmt.Println("时间戳:", header.Time)
}

func main() {
	fmt.Println("=== 测试连接 ===")
	testConnection()

	fmt.Println("\n=== 选择监听模式 ===")
	useWebSocket := true // 设置为 false 使用 HTTP 轮询

	if useWebSocket {
		fmt.Println("使用 WebSocket 实时订阅模式")
		subscribeNewHeads()
	} else {
		fmt.Println("使用 HTTP 轮询模式")
		pollNewHeads()
	}
}

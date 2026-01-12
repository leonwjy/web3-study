package lesson

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

// 查询区块
func main() {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/64820e5ca95c49c9b5022409694643ce")
	if err != nil {
		log.Fatal(err)
	}

	header, err := client.HeaderByNumber(context.Background(), nil)

	fmt.Println(header.Number.String())
	fmt.Println(header.Time)
	fmt.Println(header.Difficulty.Uint64())
	fmt.Println(header.Hash().Hex())

	if err != nil {
		log.Fatal(err)
	}

	block, err := client.BlockByNumber(context.Background(), nil)

	fmt.Println(block.Number().Uint64())
	fmt.Println(block.Time())
	fmt.Println(block.Difficulty().Uint64())
	fmt.Println(block.Hash().Hex())
	fmt.Println(block.Transactions())

	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()
}

package lesson

import (
	"log"

	"go-study/lesson/store"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	contractAddr = "0x8D4141ec2b522dE5Cf42705C3010541B4B3EC24e"
)

// 加载已部署的合约
func main() {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/64820e5ca95c49c9b5022409694643ce")
	if err != nil {
		log.Fatal(err)
	}
	storeContract, err := store.NewStore(common.HexToAddress(contractAddr), client)
	if err != nil {
		log.Fatal(err)
	}

	_ = storeContract
}

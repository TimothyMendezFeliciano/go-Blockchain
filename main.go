package main

import (
	"fmt"
	"go-blockchain/classes"
	"log"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {

	blockchainAddress := "GottaGenerateYouSoon"
	blockChain := classes.NewBlockchain(blockchainAddress)
	blockChain.Print()

	blockChain.AddTransaction("Alice", "Bob", 1.0)
	blockChain.Mining()
	blockChain.Print()

	blockChain.AddTransaction("Bob", "Charlie", 2.0)
	blockChain.AddTransaction("Charlie", "Alice", 3.56)
	blockChain.Mining()
	blockChain.Print()

	fmt.Printf("BlockchainItself %1f\n", blockChain.CalculateTotalAmount(blockchainAddress))
	fmt.Printf("Alice %1f\n", blockChain.CalculateTotalAmount("Alice"))
	fmt.Printf("Charlie %1f\n", blockChain.CalculateTotalAmount("Charlie"))
}

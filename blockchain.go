package main

import (
	"fmt"
	"log"
	"time"
)

type Block struct {
	nonce        int
	prevHash     string
	timestamp    int64
	transactions []string
}

func NewBlock(nonce int, prevHash string) *Block {

	b := new(Block)
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.prevHash = prevHash

	return b
}

func (b *Block) Print() {
	fmt.Printf("timestamp         %d\n", b.timestamp)
	fmt.Printf("nonce             %d\n", b.nonce)
	fmt.Printf("previous_hash     %s\n", b.prevHash)
	fmt.Printf("transactions      %s\n", b.transactions)
}

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	b := NewBlock(0, "init hash")
	b.Print()
}

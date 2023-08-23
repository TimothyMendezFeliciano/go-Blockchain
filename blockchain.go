package main

import (
	"fmt"
	"log"
	"strings"
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

type Blockchain struct {
	transactionPool []string
	chain           []*Block
}

func NewBlockchain() *Blockchain {
	bc := new(Blockchain)
	bc.CreateBlock(0, "init hash")
	return bc
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash string) *Block {
	b := NewBlock(nonce, previousHash)
	bc.chain = append(bc.chain, b)
	return b
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 15), i,
			strings.Repeat("=", 15))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 15))
}

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	blockChain := NewBlockchain()
	blockChain.Print()
	blockChain.CreateBlock(1, "Initial Hash 1")
	blockChain.Print()
	blockChain.CreateBlock(2, "Secondary Hash 2")
	blockChain.Print()
	blockChain.CreateBlock(3, "Third Hash 3")
	blockChain.Print()
}

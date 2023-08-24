package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type Block struct {
	timestamp    int64
	nonce        int
	prevHash     [32]byte
	transactions []*Transaction
}

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 1.0
)

func NewBlock(nonce int, prevHash [32]byte, transactions []*Transaction) *Block {

	b := new(Block)
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.prevHash = prevHash
	b.transactions = transactions

	return b
}

func (b *Block) Print() {
	fmt.Printf("timestamp         %d\n", b.timestamp)
	fmt.Printf("nonce             %d\n", b.nonce)
	fmt.Printf("previous_hash     %x\n", b.prevHash)
	for _, t := range b.transactions {
		t.Print()
	}
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256([]byte(m))
}
func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash [32]byte       `json:"previous_hash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: b.prevHash,
		Transactions: b.transactions,
	})
}

type Blockchain struct {
	transactionPool   []*Transaction
	chain             []*Block
	blockchainAddress string
}

func NewBlockchain(blockchainAddress string) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.blockchainAddress = blockchainAddress
	bc.CreateBlock(0, b.Hash())
	return bc
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

func (bc *Blockchain) AddTransaction(sender, recipient string, value float32) {
	t := NewTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, t)
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)

	for _, t := range bc.transactionPool {
		transactions = append(transactions,
			NewTransaction(t.senderAddress, t.recipientAddress, t.value))
	}

	return transactions
}

func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{0, nonce, previousHash, transactions}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())

	return guessHashStr[:difficulty] == zeros
}

func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0

	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}

	return nonce
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), i,
			strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

func (bc *Blockchain) Mining() bool {
	bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD)
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)

	log.Println("action=mining, status=success")

	return true
}

func (bc *Blockchain) CalculateTotalAmount(address string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value
			if address == t.recipientAddress {
				totalAmount += value
			}

			if address == t.senderAddress {
				totalAmount -= value
			}
		}
	}

	return totalAmount
}

type Transaction struct {
	recipientAddress string
	senderAddress    string
	value            float32
}

func NewTransaction(sender, recipient string, value float32) *Transaction {
	return &Transaction{recipientAddress: recipient, senderAddress: sender, value: value}
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("=", 40))
	fmt.Printf(" sender_blockchain_address     %s\n", t.senderAddress)
	fmt.Printf(" recipient_blockchain_address     %s\n", t.recipientAddress)
	fmt.Printf(" value     %.1f\n", t.value)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.senderAddress,
		Recipient: t.recipientAddress,
		Value:     t.value,
	})
}

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {

	blockchainAddress := "GottaGenerateYouSoon"
	blockChain := NewBlockchain(blockchainAddress)
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

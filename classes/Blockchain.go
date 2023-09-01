package classes

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"go-blockchain/constants"
	"go-blockchain/utils"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Blockchain struct {
	transactionPool   []*Transaction
	chain             []*Block
	BlockchainAddress string
	port              uint16
	mux               sync.Mutex

	neighbours    []string
	muxNeighbours sync.Mutex
}

type TransactionRequest struct {
	SenderBlockchainAddress    *string  `json:"sender_blockchain_address"`
	RecipientBlockchainAddress *string  `json:"recipient_blockchain_address"`
	SenderPublicKey            *string  `json:"sender_public_key"`
	Value                      *float32 `json:"value"`
	Signature                  *string  `json:"signature"`
}

type AmountResponse struct {
	Amount float32 `json:"amount"`
}

func (ar *AmountResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Amount float32 `json:"amount"`
	}{Amount: ar.Amount})
}

func NewBlockchain(blockchainAddress string, port uint16) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.BlockchainAddress = blockchainAddress
	bc.CreateBlock(0, b.Hash())
	bc.port = port
	return bc
}

func (bc *Blockchain) Run() {
	bc.StartSyncNeighbours()
}

func (bc *Blockchain) SetNeighbours() {
	bc.neighbours = utils.FindNeighbors(
		utils.GetHost(), bc.port,
		constants.NEIGHBOR_IP_RANGE_START, constants.NEIGHBOR_IP_RANGE_END,
		constants.BLOCKCHAIN_PORT_RANGE_START, constants.BLOCKCHAIN_PORT_RANGE_END)
	log.Printf("%v", bc.neighbours)
}

func (bc *Blockchain) SyncNeighbours() {
	bc.muxNeighbours.Lock()
	defer bc.muxNeighbours.Unlock()
	bc.SetNeighbours()
}

func (bc *Blockchain) StartSyncNeighbours() {
	bc.SyncNeighbours()
	_ = time.AfterFunc(time.Second*constants.BLOCKCHAIN_NEIGHBOUR_SYNC_TIME_SEC, bc.StartSyncNeighbours)
}

func (bc *Blockchain) TransactionPool() []*Transaction {
	return bc.transactionPool
}

func (bc *Blockchain) ClearTransactionPool() {
	bc.transactionPool = bc.transactionPool[:0]
}

func (bc *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"chains"`
	}{
		Blocks: bc.chain,
	})
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}

	for _, node := range bc.neighbours {
		endpoint := fmt.Sprintf("http://%s/transactions", node)
		client := &http.Client{}
		request, _ := http.NewRequest("DELETE", endpoint, nil)
		response, _ := client.Do(request)
		log.Printf("%v", response)
	}
	return b
}

func (bc *Blockchain) CreateTransaction(sender, recipient string, value float32,
	senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	isTransacted := bc.AddTransaction(sender, recipient, value, senderPublicKey, s)

	// TODO
	// Sync to all the servers

	if isTransacted {
		for _, n := range bc.neighbours {
			publicKeyStr := fmt.Sprintf("%064x%064x", senderPublicKey.X.Bytes(), senderPublicKey.Y.Bytes())

			signatureStr := s.String()

			blockTransaction := &TransactionRequest{&sender, &recipient, &publicKeyStr, &value, &signatureStr}
			message, _ := json.Marshal(blockTransaction)
			buffer := bytes.NewBuffer(message)

			endpoint := fmt.Sprintf("http://%s/transactions", n)
			client := &http.Client{}
			request, _ := http.NewRequest("PUT", endpoint, buffer)
			response, _ := client.Do(request)
			log.Printf("%v", response)
		}
	}
	return isTransacted
}

func (bc *Blockchain) AddTransaction(sender, recipient string, value float32,
	senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	t := NewTransaction(sender, recipient, value)

	if sender == constants.MINING_SENDER {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}

	if bc.VerifyTransactionSignature(senderPublicKey, s, t) {
		//if bc.CalculateTotalAmount(sender) < value {
		//	log.Println("ERROR: Not enough balance in a wallet.")
		//	return false
		//}
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	} else {
		log.Println("Error: Verify Transaction")
	}
	return false
}

func (bc *Blockchain) VerifyTransactionSignature(
	senderPublicKey *ecdsa.PublicKey, s *utils.Signature, t *Transaction) bool {
	marshal, _ := json.Marshal(t)
	hash := sha256.Sum256([]byte(marshal))
	return ecdsa.Verify(senderPublicKey, hash[:], s.R, s.S)
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)

	for _, t := range bc.transactionPool {
		transactions = append(transactions,
			NewTransaction(t.SenderAddress(), t.RecipientAddress(), t.Value()))
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

	for !bc.ValidProof(nonce, previousHash, transactions, constants.MINING_DIFFICULTY) {
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
	bc.mux.Lock()
	defer bc.mux.Unlock()

	if len(bc.transactionPool) == 0 {
		return false
	}
	bc.AddTransaction(constants.MINING_SENDER, bc.BlockchainAddress, constants.MINING_REWARD, nil, nil)
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)

	log.Println("action=mining, status=success")

	return true
}

func (bc *Blockchain) StartMining() {
	bc.Mining()
	_ = time.AfterFunc(time.Second*constants.MINING_TIMER_SEC, bc.StartMining)
}

func (bc *Blockchain) CalculateTotalAmount(address string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.Value()
			if address == t.RecipientAddress() {
				totalAmount += value
			}

			if address == t.SenderAddress() {
				totalAmount -= value
			}
		}
	}

	return totalAmount
}

func (tr *TransactionRequest) Validate() bool {
	return tr.RecipientBlockchainAddress != nil ||
		tr.SenderPublicKey != nil ||
		tr.Signature != nil ||
		tr.Value != nil ||
		tr.SenderBlockchainAddress != nil
}

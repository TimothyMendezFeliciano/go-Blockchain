package classes

type Block struct {
	timestamp    int64
	nonce        int
	prevHash     [32]byte
	transactions []*Transaction
}

type Blockchain struct {
	transactionPool   []*Transaction
	chain             []*Block
	BlockchainAddress string
}

type Transaction struct {
	recipientAddress string
	senderAddress    string
	value            float32
}

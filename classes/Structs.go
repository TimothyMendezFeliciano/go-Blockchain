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
	blockchainAddress string
}

type Transaction struct {
	RecipientAddress string
	SenderAddress    string
	Value            float32
}

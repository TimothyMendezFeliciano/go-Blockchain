package main

import (
	"encoding/json"
	"go-blockchain/classes"
	"go-blockchain/utils"
	"go-blockchain/wallet"
	"io"
	"log"
	"net/http"
	"strconv"
)

var cache map[string]*classes.Blockchain = make(map[string]*classes.Blockchain)

type BlockchainServer struct {
	port uint16
}

func NewBlockchainServer(port uint16) *BlockchainServer {
	return &BlockchainServer{port}
}

func (bcs *BlockchainServer) Port() uint16 {
	return bcs.port
}

func (bcs *BlockchainServer) GetBlockchain() *classes.Blockchain {
	bc, ok := cache["Blockchain"]
	if !ok {
		minersWallet := wallet.NewWallet()
		bc = classes.NewBlockchain(minersWallet.BlockchainAddress(), bcs.Port())
		cache["Blockchain"] = bc
		log.Printf("private_key %v", minersWallet.PrivateKeyStr())
		log.Printf("public_key %v", minersWallet.PublicKeyStr())
		log.Printf("blockchain_address %v", minersWallet.BlockchainAddress())
	}
	return bc
}

func (bcs *BlockchainServer) GetChain(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcs.GetBlockchain()
		m, _ := bc.MarshalJSON()
		io.WriteString(w, string(m[:]))
	default:
		log.Printf("Error: Invalid HTTP Method")
	}
}

func (bcs *BlockchainServer) Transactions(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcs.GetBlockchain()
		transactions := bc.TransactionPool()
		m, _ := json.Marshal(struct {
			Transactions []*classes.Transaction `json:"transactions"`
			Length       int                    `json:"length"`
		}{
			Transactions: transactions,
			Length:       len(transactions),
		})

		io.WriteString(w, string(m))
	case http.MethodPost:
		decoder := json.NewDecoder(req.Body)
		var transaction classes.TransactionRequest
		error := decoder.Decode(&transaction)
		if error != nil {
			log.Printf("Error: %v", error)
			io.WriteString(w, "Failed to Call Transactions")
			return
		}

		if !transaction.Validate() {
			log.Printf("Missing fields")
			io.WriteString(w, "Failed to Call Transactions")
			return
		}

		publicKey := utils.PublicKeyFromString(*transaction.SenderPublicKey)
		signature := utils.SignatureFromString(*transaction.Signature)

		bc := bcs.GetBlockchain()
		isCreated := bc.CreateTransaction(*transaction.SenderBlockchainAddress, *transaction.RecipientBlockchainAddress,
			*transaction.Value, publicKey, signature)

		w.Header().Add("Content-Type", "application/json")
		var message []byte
		if !isCreated {
			w.WriteHeader(http.StatusBadRequest)
			message = utils.JsonStatus("Bad Request")
		} else {
			w.WriteHeader(http.StatusCreated)
			message = utils.JsonStatus("Transaction Succesful")
		}
		io.WriteString(w, string(message))

	default:
		log.Println("ERROR: Incorrect Methods")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) Run() {
	http.HandleFunc("/", bcs.GetChain)
	http.HandleFunc("/transactions", bcs.Transactions)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(bcs.Port())), nil))
}

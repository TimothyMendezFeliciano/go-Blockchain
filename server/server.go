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
		err := decoder.Decode(&transaction)
		if err != nil {
			log.Printf("Error: %v", err)
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
	case http.MethodPut:
		decoder := json.NewDecoder(req.Body)
		var transaction classes.TransactionRequest
		err := decoder.Decode(&transaction)
		if err != nil {
			log.Printf("Error: %v", err)
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
		isUpdated := bc.AddTransaction(*transaction.SenderBlockchainAddress, *transaction.RecipientBlockchainAddress,
			*transaction.Value, publicKey, signature)

		w.Header().Add("Content-Type", "application/json")
		var message []byte
		if !isUpdated {
			w.WriteHeader(http.StatusBadRequest)
			message = utils.JsonStatus("Bad Request")
		} else {
			w.WriteHeader(http.StatusCreated)
			message = utils.JsonStatus("Transaction Succesful")
		}
		io.WriteString(w, string(message))
	case http.MethodDelete:
		bc := bcs.GetBlockchain()
		bc.ClearTransactionPool()
		io.WriteString(w, string(utils.JsonStatus("Pool Cleared")))
	default:
		log.Println("ERROR: Incorrect Methods")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) Mine(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		bc := bcs.GetBlockchain()
		isMined := bc.Mining()
		var message []byte
		if !isMined {
			w.WriteHeader(http.StatusBadRequest)
			message = utils.JsonStatus("Failed to Mine.")
		} else {
			message = utils.JsonStatus("Mining Succesfully.")
		}

		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(message))
	default:
		log.Println("ERROR: Incorrect Methods")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) StartMine(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		bc := bcs.GetBlockchain()
		bc.StartMining()
		var message = utils.JsonStatus("Mining Success`fully")

		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(message))
	default:
		log.Println("ERROR: Incorrect Methods")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) Amount(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		blockchainAddress := req.URL.Query().Get("blockchain_address")
		bc := bcs.GetBlockchain()
		amount := bc.CalculateTotalAmount(blockchainAddress)
		ar := classes.AmountResponse{amount}
		m, _ := ar.MarshalJSON()

		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(m[:]))
	default:
		log.Println("ERROR: Incorrect Methods")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) Run() {
	bcs.GetBlockchain().Run()
	http.HandleFunc("/", bcs.GetChain)
	http.HandleFunc("/transactions", bcs.Transactions)
	http.HandleFunc("/mine", bcs.Mine)
	http.HandleFunc("/mine/start", bcs.StartMine)
	http.HandleFunc("/amount", bcs.Amount)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(bcs.Port())), nil))
}

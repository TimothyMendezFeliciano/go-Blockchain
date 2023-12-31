package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-blockchain/classes"
	"go-blockchain/utils"
	"go-blockchain/wallet"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
)

const tempDir = "wallet_server/templates"

type WalletServer struct {
	port    uint16
	gateway string
}

func NewWalletServer(port uint16, gateway string) *WalletServer {
	return &WalletServer{port, gateway}
}

func (ws *WalletServer) Port() uint16 {
	return ws.port
}

func (ws *WalletServer) Gateway() string {
	return ws.gateway
}

func (ws *WalletServer) Index(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		t, _ := template.ParseFiles(path.Join(tempDir, "index.html"))
		t.Execute(w, "")

	default:
		log.Printf("Error: Invalid HTTP Method")
	}
}

func (ws *WalletServer) Wallet(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")
		newWallet := wallet.NewWallet()
		m, _ := newWallet.MarshalJSON()
		io.WriteString(w, string(m[:]))

	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) CreateTransaction(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(req.Body)
		var t wallet.TransactionRequest
		error := decoder.Decode(&t)

		if error != nil {
			log.Printf("ERROR: decoded value %v", error)
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		if !t.Validate() {
			log.Printf("ERROR: Invalid value %v", t)
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}

		publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
		privateKey := utils.PrivateKeyFromString(*t.SenderPrivateKey, publicKey)
		value, err := strconv.ParseFloat(*t.Value, 32)
		if err != nil {
			log.Println("Error: Parsing Error")
			io.WriteString(w, string(utils.JsonStatus("Failed to Parse Transaction")))
			return
		}

		value32 := float32(value)

		w.Header().Add("Content-Type", "application/json")

		transaction := wallet.NewTransaction(*t.SenderBlockchainAddress, *t.RecipientBlockchainAddress, value32, privateKey, publicKey)
		signature := transaction.GenerateSignature()
		signatureString := signature.String()

		blockTransaction := &classes.TransactionRequest{
			t.SenderBlockchainAddress,
			t.RecipientBlockchainAddress,
			t.SenderPublicKey,
			&value32,
			&signatureString,
		}

		m, _ := json.Marshal(blockTransaction)

		buff := bytes.NewBuffer(m)
		response, _ := http.Post(ws.Gateway()+"/transactions", "application/json", buff)

		if response.StatusCode == 201 {
			io.WriteString(w, string(utils.JsonStatus("success")))
			return
		}
		io.WriteString(w, string(utils.JsonStatus("fail")))

	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error: Invalid HTTP Method")
	}
}

func (ws *WalletServer) WalletAmount(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		blockchainAddress := req.URL.Query().Get("blockchain_address")
		endpoint := fmt.Sprintf("%s/amount", ws.Gateway())

		client := &http.Client{}
		bcsRequest, _ := http.NewRequest("GET", endpoint, nil)
		query := bcsRequest.URL.Query()
		query.Add("blockchain_address", blockchainAddress)
		bcsRequest.URL.RawQuery = query.Encode()

		bcsResponse, error := client.Do(bcsRequest)

		if error != nil {
			log.Printf("Error: %v", error)
			io.WriteString(w, string(utils.JsonStatus("Failed to Get Wallet Amount")))
			return
		}

		w.Header().Add("Content-Type", "application/json")

		if bcsResponse.StatusCode == 200 {
			decoder := json.NewDecoder(bcsResponse.Body)
			var bar classes.AmountResponse
			error := decoder.Decode(&bar)
			if error != nil {
				log.Printf("Error: %v", error)
				io.WriteString(w, string(utils.JsonStatus("Failed to Get Wallet Amount")))
				return
			}

			m, _ := json.Marshal(struct {
				Message string  `json:"message"`
				Amount  float32 `json:"amount"`
			}{
				Message: "Success",
				Amount:  bar.Amount,
			})
			io.WriteString(w, string(m[:]))
		} else {
			io.WriteString(w, string(utils.JsonStatus("Failure")))
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error: Invalid HTTP Method")
	}
}

func (ws *WalletServer) Run() {
	http.HandleFunc("/", ws.Index)
	http.HandleFunc("/wallet", ws.Wallet)
	http.HandleFunc("/wallet/amount", ws.WalletAmount)
	http.HandleFunc("/transactions", ws.CreateTransaction)
	log.Println(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ws.Port())), nil))
}

package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"go-blockchain/utils"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	blockchainAddress string
}

type TransactionRequest struct {
	SenderPrivateKey           *string `json:"sender_private_key"`
	SenderPublicKey            *string `json:"sender_public_key"`
	SenderBlockchainAddress    *string `json:"sender_blockchain_address"`
	RecipientBlockchainAddress *string `json:"recipient_blockchain_address"`
	Value                      *string `json:"value"`
}

func NewWallet() *Wallet {
	w := new(Wallet)
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	w.privateKey = privateKey
	w.publicKey = &w.privateKey.PublicKey

	hash2 := sha256.New()
	hash2.Write(w.publicKey.X.Bytes())
	hash2.Write(w.publicKey.Y.Bytes())
	digest2 := hash2.Sum(nil)

	hash3 := ripemd160.New()
	hash3.Write(digest2)
	digest3 := hash3.Sum(nil)

	versionD4 := make([]byte, 21)
	versionD4[0] = 0x00
	copy(versionD4[1:], digest3[:])

	hash5 := sha256.New()
	hash5.Write(versionD4)
	digest5 := hash5.Sum(nil)

	hash6 := sha256.New()
	hash6.Write(digest5)
	digest6 := hash6.Sum(nil)

	checksum := digest6[:4]
	dc8 := make([]byte, 25)
	copy(dc8[:21], versionD4[:])
	copy(dc8[21:], checksum[:])

	address := base58.Encode(dc8)
	w.blockchainAddress = address
	return w
}

func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%064x%064x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}

func (w *Wallet) BlockchainAddress() string {
	return w.blockchainAddress
}

type Transaction struct {
	recipientAddress string
	senderAddress    string
	value            float32
	senderPrivateKey *ecdsa.PrivateKey
	senderPublicKey  *ecdsa.PublicKey
}

func NewTransaction(senderAddress, recipientAddress string, value float32, senderPrivateKey *ecdsa.PrivateKey, senderPublicKey *ecdsa.PublicKey) *Transaction {
	return &Transaction{recipientAddress, senderAddress, value, senderPrivateKey, senderPublicKey}
}

func (t *Transaction) RecipientAddress() string {
	return t.recipientAddress
}

func (t *Transaction) GenerateSignature() *utils.Signature {
	marshal, _ := json.Marshal(t)
	hash := sha256.Sum256([]byte(marshal))
	r, s, _ := ecdsa.Sign(rand.Reader, t.senderPrivateKey, hash[:])

	return &utils.Signature{R: r, S: s}
}

func (w *Wallet) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PrivateKey        string `json:"private_key"`
		PublicKey         string `json:"public_key"`
		BlockchainAddress string `json:"blockchain_address"`
	}{
		PrivateKey:        w.PublicKeyStr(),
		PublicKey:         w.PublicKeyStr(),
		BlockchainAddress: w.BlockchainAddress(),
	})
}

func (tr *TransactionRequest) Validate() bool {
	if tr.RecipientBlockchainAddress == nil ||
		tr.SenderBlockchainAddress == nil ||
		tr.SenderPrivateKey == nil ||
		tr.SenderPublicKey == nil ||
		tr.Value == nil {
		return false
	} else {
		return true
	}
}

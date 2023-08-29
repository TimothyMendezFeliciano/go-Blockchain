package classes

import (
	"encoding/json"
	"fmt"
	"strings"
)

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

func (t *Transaction) RecipientAddress() string {
	return t.recipientAddress
}
func (t *Transaction) SenderAddress() string {
	return t.senderAddress
}
func (t *Transaction) Value() float32 {
	return t.value
}

package classes

import (
	"encoding/json"
	"fmt"
	"strings"
)

func NewTransaction(sender, recipient string, value float32) *Transaction {
	return &Transaction{RecipientAddress: recipient, SenderAddress: sender, Value: value}
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("=", 40))
	fmt.Printf(" sender_blockchain_address     %s\n", t.SenderAddress)
	fmt.Printf(" recipient_blockchain_address     %s\n", t.RecipientAddress)
	fmt.Printf(" value     %.1f\n", t.Value)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.SenderAddress,
		Recipient: t.RecipientAddress,
		Value:     t.Value,
	})
}

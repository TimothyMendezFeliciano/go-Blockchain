package _test

import (
	"go-blockchain/classes"
	"testing"
)

func TestNewTransaction(t *testing.T) {
	senderAddress := "Alice"
	recipientAddress := "Bob"
	var value float32 = 12.34
	expected := &classes.Transaction{
		RecipientAddress: recipientAddress,
		SenderAddress:    senderAddress,
		Value:            value,
	}
	result := classes.NewTransaction(senderAddress, recipientAddress, value)

	if result.SenderAddress != expected.SenderAddress &&
		result.RecipientAddress != expected.RecipientAddress &&
		result.Value != expected.Value {
		t.Errorf("\"NewTransaction('%s')\" FAILED, expected -> %v, got -> %v", "_", expected, result)
	} else {
		t.Logf("\"NewTransaction('%s')\" SUCCEDED, expected -> %v, got -> %v", "_", expected, result)
	}
}

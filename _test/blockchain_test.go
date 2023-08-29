package _test

import (
	"go-blockchain/classes"
	"go-blockchain/wallet"
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

func TestNewBlockchain(t *testing.T) {
	blockchainAddress := "konosuba"
	nonce := 12

	b := &classes.Block{}

	bcExpected := new(classes.Blockchain)
	bcExpected.BlockchainAddress = blockchainAddress
	bcExpected.CreateBlock(nonce, b.Hash())

	result := classes.NewBlockchain(blockchainAddress)

	if result.BlockchainAddress == bcExpected.BlockchainAddress {
		t.Logf("\"BlockchainAddress('%s')\" SUCCEDED, expected -> %v, got -> %v", "_", bcExpected, result)
	}
}

func TestNewWallet(t *testing.T) {
	w := wallet.NewWallet()

	if w.PublicKeyStr() != w.PrivateKeyStr() {
		t.Logf("\"BlockchainWallet\" CREATED, PrivateKey -> %v, PublicKey -> %v", w.PrivateKeyStr(), w.PublicKeyStr())
	} else {
		t.Errorf("Wallet Created Incorrectly")
	}
}

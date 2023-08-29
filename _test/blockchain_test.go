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
	result := classes.NewTransaction(senderAddress, recipientAddress, value)

	if result.SenderAddress() != senderAddress &&
		result.RecipientAddress() != recipientAddress &&
		result.Value() != value {
		t.Errorf("\"NewTransaction('%s')\" FAILED, expected -> %v, got -> %v", "_", classes.NewTransaction(senderAddress, recipientAddress, value), result)
	} else {
		t.Logf("\"NewTransaction('%s')\" SUCCEDED, expected -> %v, got -> %v", "_", classes.NewTransaction(senderAddress, recipientAddress, value), result)
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

	if w.PublicKeyStr() != w.PrivateKeyStr() && len(w.BlockchainAddress()) >= 1 {
		t.Logf("\"BlockchainWallet\" CREATED, PrivateKey -> %v, PublicKey -> %v, Address -> %v", w.PrivateKeyStr(), w.PublicKeyStr(), w.BlockchainAddress())
	} else {
		t.Errorf("Wallet Created Incorrectly")
	}
}

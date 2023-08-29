package _test

import (
	"go-blockchain/classes"
	"go-blockchain/wallet"
	"testing"
)

///

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

func TestNewTransaction(t *testing.T) {
	w := wallet.NewWallet()

	transaction := w.NewTransaction(w.BlockchainAddress(), "Bob", 1.0, w.PrivateKey(), w.PublicKey())
	signature := transaction.GenerateSignature()

	if len(signature.String()) >= 1 {
		t.Logf("\"Transaction\" SUCCESFUL, Signature -> %v", signature)
	} else {
		t.Errorf("\"Transaction\" ERROR, Signature -> %v", signature)
	}

	if transaction.RecipientAddress() == "Bob" {
		t.Logf("\"Transaction\" SUCCESFUL, RecipientAddress -> %v", transaction.RecipientAddress())
	} else {
		t.Errorf("\"Transaction\" ERROR, Signature -> %v", signature)
	}
}

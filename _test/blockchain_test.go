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

	transaction := wallet.NewTransaction(w.BlockchainAddress(), "Bob", 1.0, w.PrivateKey(), w.PublicKey())
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

func TestAddTransaction(t *testing.T) {
	AliceWallet := wallet.NewWallet()
	BobWallet := wallet.NewWallet()
	MinerWallet := wallet.NewWallet()

	blockchain := classes.NewBlockchain(MinerWallet.BlockchainAddress())
	blockchain.Mining()

	AliceBeforeTransactionAmount := blockchain.CalculateTotalAmount(AliceWallet.BlockchainAddress())
	BobBeforeTransactionAmount := blockchain.CalculateTotalAmount(BobWallet.BlockchainAddress())
	transaction := wallet.NewTransaction(AliceWallet.BlockchainAddress(), BobWallet.BlockchainAddress(), 1.0, AliceWallet.PrivateKey(), AliceWallet.PublicKey())

	isAdded := blockchain.AddTransaction(AliceWallet.BlockchainAddress(), BobWallet.BlockchainAddress(), 1.0, AliceWallet.PublicKey(), transaction.GenerateSignature())

	if isAdded {
		t.Logf("\"Add Transaction\" SUCCESSFUL, Sender -> %v, Recipient -> %v, transaction -> %v", AliceWallet.BlockchainAddress(), BobWallet.BlockchainAddress(), transaction)
	} else {
		t.Errorf("\"Add Transaction\" ERROR, Sender -> %v, Recipient -> %v, transaction -> %v", AliceWallet.BlockchainAddress(), BobWallet.BlockchainAddress(), transaction)
	}

	AliceAfterTransactionAmount := blockchain.CalculateTotalAmount(AliceWallet.BlockchainAddress())
	BobAfterTransactionAmount := blockchain.CalculateTotalAmount(BobWallet.BlockchainAddress())
	if AliceBeforeTransactionAmount >= AliceAfterTransactionAmount &&
		BobBeforeTransactionAmount <= BobAfterTransactionAmount {
		t.Logf("\"Transaction Amount\" SUCCESSUL, Address -> %v, Before -> %v, After -> %v", AliceWallet.BlockchainAddress(), AliceBeforeTransactionAmount, AliceAfterTransactionAmount)
		t.Logf("\"Transaction Amount\" SUCCESSUL, Address -> %v, Before -> %v, After -> %v", BobWallet.BlockchainAddress(), BobBeforeTransactionAmount, BobAfterTransactionAmount)
	} else {
		t.Errorf("\"Transaction Amount\" ERROR, Address -> %v, Before -> %v, After -> %v", AliceWallet.BlockchainAddress(), AliceBeforeTransactionAmount, AliceAfterTransactionAmount)
		t.Errorf("\"Transaction Amount\" ERROR, Address -> %v, Before -> %v, After -> %v", BobWallet.BlockchainAddress(), BobBeforeTransactionAmount, BobAfterTransactionAmount)
	}
}

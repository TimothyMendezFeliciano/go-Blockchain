package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"
)

type Signature struct {
	R *big.Int
	S *big.Int
}

func (s *Signature) String() string {
	return fmt.Sprintf("%064x%064x", s.R, s.S)
}

func String2BigIntTuple(s string) (big.Int, big.Int) {
	bigX, _ := hex.DecodeString(s[:64])
	bigY, _ := hex.DecodeString(s[64:])

	var bigIntX big.Int
	var bigIntY big.Int

	_ = bigIntX.SetBytes(bigX)
	_ = bigIntY.SetBytes(bigY)

	return bigIntX, bigIntY
}

func SignatureFromString(s string) *Signature {
	x, y := String2BigIntTuple(s)
	return &Signature{&x, &y}
}

func PublicKeyFromString(s string) *ecdsa.PublicKey {
	x, y := String2BigIntTuple(s)
	return &ecdsa.PublicKey{elliptic.P256(), &x, &y}
}

func PrivateKeyFromString(s string, publicKey *ecdsa.PublicKey) *ecdsa.PrivateKey {
	bytes, _ := hex.DecodeString(s[:])
	var bigInt big.Int

	_ = bigInt.SetBytes(bytes)

	return &ecdsa.PrivateKey{*publicKey, &bigInt}
}

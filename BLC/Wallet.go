package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/ripemd160"
)

// wallet management related files
// checksum length
const addressCheckSumLen = 4

// Wallet basic structure
type Wallet struct {
	// 1. Private key
	PrivateKey ecdsa.PrivateKey
	// 2. public key
	PublicKey []byte
}

// create a wallet
func NewWallt() *Wallet {
	// public-private key assignment
	privateKey, publicKey := newKeyPair()
	return &Wallet{PrivateKey: privateKey, PublicKey: publicKey}
}

// Generate public-private key pair from wallet
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	// 1. Get an ellipse
	curve := elliptic.P256()
	// 2. Generate private key through ellipse correlation algorithm
	priv, err := ecdsa.GenerateKey(curve, rand.Reader)
	if nil != err {
		log.Panicf("ecdsa generate private key failed! %v\n", err)
	}
	// 3. Generate public key from private key
	pubKey := append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)
	return *priv, pubKey
}

// generate address

// implement double hash
func Ripemd160Hash(pubKey []byte) []byte {
	// 1. sha256
	hash256 := sha256.New()
	hash256.Write(pubKey)
	hash := hash256.Sum(nil)
	// 2. ripemd160
	rmd160 := ripemd160.New()
	rmd160.Write(hash)
	return rmd160.Sum(nil)
}

// generate checksum
func CheckSum(input []byte) []byte {
	first_hash := sha256.Sum256(input)
	second_hash := sha256.Sum256(first_hash[:])
	return second_hash[:addressCheckSumLen]
}

// Get the address through the wallet (public key)
func (w *Wallet) GetAddress() []byte {
	// 1. Get hash160
	ripemd160Hash := Ripemd160Hash(w.PublicKey)
	// 2. Get the checksum
	checkSumBytes := CheckSum(ripemd160Hash)
	// 3. Address composition member splicing
	addressBytes := append(ripemd160Hash, checkSumBytes...)
	// 4. base58 encoding
	b58Bytes := Base58Encode(addressBytes)
	return b58Bytes
}

// Determine the validity of the address
func IsValidForAddress(addressBytes []byte) bool {
	// 1. The address is decoded by base58Decode (length is 24)
	pubkey_checkSumByte := Base58Decode(addressBytes)
	// 2. Split, checksum verification
	checkSumBytes := pubkey_checkSumByte[len(pubkey_checkSumByte)-addressCheckSumLen:]
	// Pass in ripemdhash160 to generate checksum
	ripemd160hash := pubkey_checkSumByte[:len(pubkey_checkSumByte)-addressCheckSumLen]
	// 3. Generate
	checkBytes := CheckSum(ripemd160hash)
	// 4. Compare
	if bytes.Compare(checkSumBytes, checkBytes) == 0 {
		return true
	}
	return false
}

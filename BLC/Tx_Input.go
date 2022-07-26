package BLC

import "bytes"

// transaction input management

// input structure
type TxInput struct {
	// transaction hash (not the current transaction hash)
	TxHash []byte
	// The output index number of the last transaction referenced
	Vout int
	// digital signature
	Signature []byte
	// public key
	PublicKey []byte
}

// Pass the hash 160 for judgment
func (in *TxInput) UnLockRipemd160Hash(ripemd160Hash []byte) bool {
	// Get the ripemd160 hash value of the input
	inputRipemd160Hash := Ripemd160Hash(in.PublicKey)
	return bytes.Compare(inputRipemd160Hash, ripemd160Hash) == 0
}

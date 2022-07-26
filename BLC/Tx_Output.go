package BLC

import "bytes"

// transaction output management

// output structure
type TxOutput struct {
	// Amount (capitalized to export amount)
	Value int

	//ScriptPubkey string
	// username (owner of the UTXO)
	Ripemd160Hash []byte
}

// Verify that the current UTXO belongs to the specified address
//func (txOutput *TxOutput) CheckPubkeyWithAddress(address string) bool {
// return address == txOutput.ScriptPubkey
//}

// output authentication
func (TxOutput *TxOutput) UnLockScriptPubkeyWithAddress(address string) bool {
	// convert
	hash160 := StringToHash160(address)
	return bytes.Compare(hash160, TxOutput.Ripemd160Hash) == 0
}

// create a new output object
func NewTxOutput(value int, address string) *TxOutput {
	txOutput := &TxOutput{}
	hash160 := StringToHash160(address)
	txOutput.Value = value
	txOutput.Ripemd160Hash = hash160
	return txOutput
}

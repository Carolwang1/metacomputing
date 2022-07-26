package BLC

// UTXO structure management
type UTXO struct {
	// Transaction hash corresponding to UTXO
	TxHash []byte
	// The index of the UTXO in the output list of the transaction to which it belongs
	Index int
	// Output itself
	Output *TxOutput
}

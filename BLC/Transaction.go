package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/labstack/gommon/log"
	"math/big"
	"time"
)

// transaction management file

// define a basic transaction structure
type Transaction struct {
	// transaction hash (identifier)
	TxHash []byte
	// input list
	Vins []*TxInput
	// output list
	Vouts []*TxOutput
}

// Implement coinbase transaction
func NewCoinbaseTransaction(address string) *Transaction {

	// input
	// coinbase features
	// txHash: nil
	// vout: -1 (in order to judge whether it is a coinbase transaction)
	// ScriptSig: system reward
	txInput := &TxInput{[]byte{}, -1, nil, nil}
	// output:
	// value:
	// address

	//txOutput := &TxOutput{10, StringToHash160(address)}
	txOutput := NewTxOutput(0, address)
	// input and output assembly transaction
	txCoinbase := &Transaction{
		nil,
		[]*TxInput{txInput},
		[]*TxOutput{txOutput},
	}
	// transaction hash generation
	txCoinbase.HashTransaction()
	return txCoinbase
}

// Generate transaction hash (transaction serialization)
// Transaction hash values ​​generated at different times are different
func (tx *Transaction) HashTransaction() {
	var result bytes.Buffer
	// set encoding object
	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(tx); err != nil {
		log.Panicf("tx Hash encoded failed %v\n", err)
	}
	// Add a timestamp, not adding it will cause all coinbase transaction hashes to be exactly the same
	tm := time.Now().UnixNano()
	// raw data used to generate hash
	txHashBytes := bytes.Join([][]byte{result.Bytes(), IntToHex(tm)}, []byte{})
	// generate hash value
	hash := sha256.Sum256(txHashBytes)
	tx.TxHash = hash[:]
}

// Generate a normal transfer transaction
func NewSimpleTransaction(from string, to string, amount int,
	bc *BlockChain, txs []*Transaction, nodeID string) *Transaction {
	var txInputs []*TxInput   // input list
	var txOutputs []*TxOutput // output list
	// Call the costable UTXO function
	money, spendableUTXODic := bc.FindSpendableUTXO(from, amount, txs)
	fmt.Printf("money : %v\n", money)
	// Get the wallet collection object
	wallets := NewWallets(nodeID)
	// Find the corresponding wallet structure
	wallet := wallets.Wallets[from]
	// input
	for txHash, indexArray := range spendableUTXODic {
		txHashesBytes, err := hex.DecodeString(txHash)
		if nil != err {
			log.Panicf("decode string to []byte failed! %v\n", err)
		}
		// loop through the index list
		for _, index := range indexArray {
			txInput := &TxInput{txHashesBytes, index, nil, wallet.PublicKey}
			txInputs = append(txInputs, txInput)
		}
	}
	// output (transfer source)
	//txOutput := &TxOutput{amount, to}
	txOutput := NewTxOutput(amount, to)
	txOutputs = append(txOutputs, txOutput)
	// output (change)
	if money > amount {
		//txOutput = &TxOutput{money - amount, from}
		txOutput = NewTxOutput(money-amount, from)
		txOutputs = append(txOutputs, txOutput)
	} else if money == amount {

	} else {
		log.Panicf("Insufficient balance...\n")
	}

	tx := Transaction{nil, txInputs, txOutputs}
	tx.HashTransaction() // Generate a complete transaction

	// sign the transaction
	bc.SignTransaction(&tx, wallet.PrivateKey)
	return &tx
}

// Determine if the specified transaction is a coinbase transaction
func (tx *Transaction) IsCoinbaseTransaction() bool {
	return tx.Vins[0].Vout == -1 && len(tx.Vins[0].TxHash) == 0
}

// transaction signature
// prevTxs : represents the transaction to which all OUTPUTs referenced by the input of the current transaction belong
func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey,
	prevTxs map[string]Transaction) {
	// Process the input to ensure the correctness of the transaction
	// Check if the transaction hash referenced by each input in tx is included in prevTxsa
	// If it is not included, it means that the transaction has been modified by someone
	for _, vin := range tx.Vins {
		if prevTxs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panicf("ERROR:Prev transaction is not correct!\n")
		}
	}
	// extract properties that need to be signed
	txCopy := tx.TrimmedCopy()
	// handle the input of the transaction copy
	for vin_id, vin := range txCopy.Vins {
		// get related transactions
		prevTx := prevTxs[hex.EncodeToString(vin.TxHash)]
		// find sender (hash of current input reference -- hash of output)
		txCopy.Vins[vin_id].PublicKey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		// Generate hash of transaction copy
		txCopy.TxHash = txCopy.Hash()
		// call the core signature function
		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, txCopy.TxHash)
		if nil != err {
			log.Panicf("sign to transaction [%x] failed! %v\n", err)
		}

		// Make up the transaction signature
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vins[vin_id].Signature = signature
	}

}

// Transaction copy, generate a copy dedicated to transaction signing
func (tx *Transaction) TrimmedCopy() Transaction {
	// Reassemble to generate a new transaction
	var inputs []*TxInput
	var outputs []*TxOutput
	// assemble input
	for _, vin := range tx.Vins {
		inputs = append(inputs, &TxInput{vin.TxHash, vin.Vout,
			nil, nil})
	}
	// assemble output
	for _, vout := range tx.Vouts {
		outputs = append(outputs, &TxOutput{vout.Value, vout.Ripemd160Hash})
	}
	txCopy := Transaction{tx.TxHash, inputs, outputs}
	return txCopy
}

// Set the hash of the transaction used for signing
func (tx *Transaction) Hash() []byte {
	txCopy := tx
	txCopy.TxHash = []byte{}
	hash := sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

// transaction serialization
func (tx *Transaction) Serialize() []byte {
	var buffer bytes.Buffer
	// create a new encoding object
	encoder := gob.NewEncoder(&buffer)
	// encode (serialize)
	if err := encoder.Encode(tx); nil != err {
		log.Panicf("serialize the tx to []byte failed! %v\n", err)
	}
	return buffer.Bytes()
}

// verify signature
func (tx *Transaction) Verify(prevTxs map[string]Transaction) bool {

	// Check if the transaction hash can be found
	for _, vin := range tx.Vins {
		if prevTxs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panicf("VERIFY ERROR : transaction verify failed!\n")
		}
	}

	// extract the same transaction signature property
	txCopy := tx.TrimmedCopy()
	// use the same ellipse
	curve := elliptic.P256()

	// Traverse the tx input and verify the output referenced by each input
	for vinId, vin := range tx.Vins {
		// get related transactions
		prevTx := prevTxs[hex.EncodeToString(vin.TxHash)]
		// find sender (hash of current input reference -- hash of output)
		txCopy.Vins[vinId].PublicKey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		// The transaction hash generated by the data to be verified must be exactly the same as the data at the time of signing
		txCopy.TxHash = txCopy.Hash()
		// In bits, the signature is a value pair, r, s represents the signature
		// So get it from the input signature
		// Get r, s. r, s have the same length
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])
		// get the public key
		// The public key consists of X, Y coordinates
		x := big.Int{}
		y := big.Int{}
		pubKeyLen := len(vin.PublicKey)
		x.SetBytes(vin.PublicKey[:(pubKeyLen / 2)])
		y.SetBytes(vin.PublicKey[(pubKeyLen / 2):])
		rawPublicKey := ecdsa.PublicKey{curve, &x, &y}
		if !ecdsa.Verify(&rawPublicKey, txCopy.TxHash, &r, &s) {
			return false
		}
	}
	return true
}

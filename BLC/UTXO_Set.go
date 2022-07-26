package BLC

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/labstack/gommon/log"
)

// UTXO persistence related management

// Bucket for depositing to utxo
const utxoTableName = "utxoTable"

// utxoSet structure (save all UTXOs in the specified blockchain)
type UTXOSet struct {
	Blockchain *BlockChain
}

// output set serialization
func (txOutputs *TXOutputs) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(txOutputs); nil != err {
		log.Panicf("serialize the utxo failed! %v\n", err)
	}

	return result.Bytes()
}

// output set return serialization
func DeserializeTXOutputs(txOutputsBytes []byte) *TXOutputs {
	var txOutputs TXOutputs
	decoder := gob.NewDecoder(bytes.NewReader(txOutputsBytes))
	if err := decoder.Decode(&txOutputs); nil != err {
		log.Panicf("deserialize the struct utxo failed! %v\n", err)
	}
	return &txOutputs
}

// reset
func (utxoSet *UTXOSet) ResetUTXOSet() {
	// Update the utxo table when it is first created
	utxoSet.Blockchain.DB.Update(func(tx *bolt.Tx) error {
		// find utxo table
		b := tx.Bucket([]byte(utxoTableName))
		if nil != b {
			err := tx.DeleteBucket([]byte(utxoTableName))
			if nil != err {
				log.Panicf("delete the utxo table failed! %v\n", err)
			}
		}

		// create
		bucket, err := tx.CreateBucket([]byte(utxoTableName))
		if nil != err {
			log.Panicf("create bucket failed! %v\n", err)
		}
		if nil != bucket {
			// Find all current UTXOs
			txOutputMap := utxoSet.Blockchain.FindUTXOMap()

			for keyHash, outputs := range txOutputMap {
				// store all UTXOs in
				txHash, _ := hex.DecodeString(keyHash)
				fmt.Printf("txHash : %x\n", txHash)

				// store in utxo table
				err := bucket.Put(txHash, outputs.Serialize())
				if nil != err {
					log.Panicf("put the utxo into table failed! %v\n", err)
				}
			}
		}
		return nil
	})
}

// Check balances
func (utxoSet *UTXOSet) GetBalance(address string) int {
	UTXOS := utxoSet.FindUTXOWithAddress(address)
	var amount int
	for _, utxo := range UTXOS {
		fmt.Printf("utxo-txhash:%x\n", utxo.TxHash)
		fmt.Printf("utxo-index:%x\n", utxo.Index)
		fmt.Printf("utxo-Ripemd160Hash:%x\n", utxo.Output.Ripemd160Hash)
		fmt.Printf("utxo-Value:%x\n", utxo.Output.Value)
		amount += utxo.Output.Value
	}
	return amount
}

// find
func (utxoSet *UTXOSet) FindUTXOWithAddress(address string) []*UTXO {
	var utxos []*UTXO
	err := utxoSet.Blockchain.DB.View(func(tx *bolt.Tx) error {
		// 1. Get the utxotable table
		b := tx.Bucket([]byte(utxoTableName))
		if nil != b {
			// cursor--cursor
			c := b.Cursor()
			// Traverse the data in the boltdb database through the cursor
			for k, v := c.First(); k != nil; k, v = c.Next() {
				txOutputs := DeserializeTXOutputs(v)
				for _, utxo := range txOutputs.TXOutputs {
					if utxo.UnLockScriptPubkeyWithAddress(address) {
						utxo_signle := UTXO{Output: utxo, TxHash: k}
						utxos = append(utxos, &utxo_signle)
					}
				}
			}
		}
		return nil
	})
	if nil != err {
		log.Panicf("find the utxo of [%s] failed! %v\n", address, err)
	}
	return utxos
}

// renew
func (utxoSet *UTXOSet) update(block *Block) {
	// get the latest block
	latest_block := block
	//latest_block := utxoSet.Blockchain.Iterator().Next()
	utxoSet.Blockchain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		c := tx.Bucket([]byte(blockTableName))

		if nil != b {
			utxo := b.Get(latest_block.Hash)
			if nil != utxo {
				return nil
			}
			// just look up the transaction list for the latest block, because every block is on the chain
			// The utxo table is updated once, so just look for the transaction in the most recent block
			for _, tx := range latest_block.Txs {
				if !tx.IsCoinbaseTransaction() {
					// 2. Delete the UTXO already referenced by the input of the current transaction
					for _, vin := range tx.Vins {
						// output that needs to be updated
						updatedOutputs := TXOutputs{}
						// Get the output of the transaction hash referenced by the specified input
						outputBytes := b.Get(vin.TxHash)
						if nil == outputBytes {
							blockBytes := c.Get(vin.TxHash)
							rawBlock := DeserializeBlock(blockBytes)
							utxoSet.update(rawBlock)
							outputBytes = b.Get(vin.TxHash)
						}
						// output list
						outs := DeserializeTXOutputs(outputBytes)
						for outIdx, out := range outs.TXOutputs {
							//if vin.Vout != outIdx {
							// updatedOutputs.TXOutputs = append(updatedOutputs.TXOutputs, out)
							//}
							if vin.Vout == outIdx {
								out.Value = 0
							}
							updatedOutputs.TXOutputs = append(updatedOutputs.TXOutputs, out)
						}
						canDel := true
						for _, output := range updatedOutputs.TXOutputs {
							if output.Value != 0 {
								canDel = false
							}
						}
						//If there is no UTXO in the transaction, delete the transaction
						if canDel {
							b.Delete(vin.TxHash)
						} else {
							// Store the updated utxo data in the database
							b.Put(vin.TxHash, updatedOutputs.Serialize())
						}
					}
				}

				// Get the newly generated transaction output in the current block
				// 1. Insert the UTXO in the latest block
				newOutputs := TXOutputs{}
				newOutputs.TXOutputs = append(newOutputs.TXOutputs, tx.Vouts...)
				b.Put(tx.TxHash, newOutputs.Serialize())
			}

		}
		return nil
	})
}

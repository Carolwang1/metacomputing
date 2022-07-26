package BLC

import (
	"bytes"
	"encoding/gob"
	"github.com/labstack/gommon/log"
	"time"
)

// Block basic structure and function management file
// Implement a basic block structure
type Block struct {
	TimeStamp     int64          //Block timestamp, representing the block time
	Hash          []byte         // current block hash
	PrevBlockHash []byte         // previous block hash
	Heigth        int64          // block height
	Txs           []*Transaction // Transaction data (transaction list)
	Nonce         int64          // The hash change value generated when running pow also represents the data dynamically modified when pow is running
}

// new block
func NewBlock(height int64, prevBlockHash []byte, txs []*Transaction) *Block {
	var block Block
	block = Block{
		TimeStamp:     time.Now().Unix(),
		Hash:          nil,
		PrevBlockHash: prevBlockHash,
		Heigth:        height,
		Txs:           txs,
	}
	// replace setHash
	// Generate a new hash via POW
	pow := NewProofOfWork(&block)
	// 执行工作量证明算法
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = int64(nonce)
	return &block
}

// Generate genesis block
func CreateGenesisBlock(txs []*Transaction) *Block {
	return NewBlock(1, nil, txs)
}

// Block structure serialization
func (block *Block) Serialize() []byte {
	var buffer bytes.Buffer
	// Create a new encoding object
	encoder := gob.NewEncoder(&buffer)
	// encoding (serialization)
	if err := encoder.Encode(block); nil != err {
		log.Panicf("serialize the block to []byte failed! %v\n", err)
	}
	return buffer.Bytes()
}

// Block data deserialization
func DeserializeBlock(blockBytes []byte) *Block {
	var block Block
	// New decoder object
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	if err := decoder.Decode(&block); nil != err {
		log.Panicf("deserialize the []byte to block failed! %v\n", err)
	}
	return &block
}

// Serialize all transaction structures in the specified block (Merkle-like hash calculation method)
func (block *Block) HashTransaction() []byte {
	var txHashes [][]byte
	//Concatenate all transaction hashes in the specified block
	for _, tx := range block.Txs {
		txHashes = append(txHashes, tx.TxHash)
	}
	// Store the transaction data in the Merkle tree, and then generate the Merkle root node
	mtree := NewMerkleTree(txHashes)
	return mtree.RootNode.Data
}

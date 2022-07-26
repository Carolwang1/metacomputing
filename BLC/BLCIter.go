package BLC

import (
	"github.com/boltdb/bolt"
	"github.com/labstack/gommon/log"
)

// Blockchain Iterator to manage files

// Iterator basic structure
type BlockChainIterator struct {
	DB          *bolt.DB //Iteration target
	CurrentHash []byte   // the hash of the current iteration target
}

// Create an iterator object
func (blc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{blc.DB, blc.Tip}
}

// Implement the iterative function next to get each block
func (bcit *BlockChainIterator) Next() *Block {
	var block *Block

	err := bcit.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if nil != b {
			currentBlockBytes := b.Get(bcit.CurrentHash)
			block = DeserializeBlock(currentBlockBytes)
			// Update the hash of the block in the iterator
			bcit.CurrentHash = block.PrevBlockHash
		}
		return nil
	})
	if nil != err {
		log.Panicf("iterator the db failed! %v\n", err)
	}
	return block
}

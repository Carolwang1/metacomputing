package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// Consensus algorithm management file

// Implement POW instances and related functions

// target difficulty value
const targetBit = 16

// structure of proof of work
type ProofOfWork struct {
	// Blocks that need consensus verification
	Block *Block
	// hash of target difficulty (big data storage)
	target *big.Int
}

// create a POW object
func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	// data length is 8 bits
	// Requirements: The first two digits need to be 0 to solve the problem
	// 1 * 2 << (8-2) = 64
	// 0100 0000
	// 0011 1111 = 63
	// 32 * 8
	// Set the target difficulty value (if the first n bits are 0, move left by 256-n bits, as long as the generated hash value is less than this 2^(256-n), it must be less than the target difficulty value)
	target = target.Lsh(target, 256-targetBit)
	return &ProofOfWork{Block: block, target: target}
}

// Execute pow, compare hashes
// Return the hash value, and the number of collisions
func (proofOfWork *ProofOfWork) Run() ([]byte, int) {
	// number of collisions
	var nonce = 0
	var hashInt big.Int
	var hash [32]byte // Generated hash value
	// Infinite loop, generating eligible hash values
	for {
		// generate preparation data
		dataBytes := proofOfWork.prepareData(int64(nonce))
		hash = sha256.Sum256(dataBytes)
		hashInt.SetBytes(hash[:])
		// Check whether the generated hash value meets the conditions
		if proofOfWork.target.Cmp(&hashInt) == 1 {
			// Found a matching hash value, break the loop
			break
		}
		nonce++
	}
	fmt.Printf("\nNumber of collisions:%d\n", nonce)
	return hash[:], nonce
}

// generate preparation data
func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	var data []byte
	// Concatenate block attributes and perform hash calculation
	timeStampBytes := IntToHex(pow.Block.TimeStamp)
	heigthBytes := IntToHex(pow.Block.Heigth)
	data = bytes.Join([][]byte{
		heigthBytes,
		timeStampBytes,
		pow.Block.PrevBlockHash,
		pow.Block.HashTransaction(),
		IntToHex(nonce),
		IntToHex(targetBit),
	}, []byte{})
	return data
}

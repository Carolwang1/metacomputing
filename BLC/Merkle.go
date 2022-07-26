package BLC

import "crypto/sha256"

// Merkle tree implementation management

type MerkleTree struct {
	// root node
	RootNode *MerkleNode
}

// merkle node structure
type MerkleNode struct {
	// left child node
	Left *MerkleNode
	// right child node
	Right *MerkleNode
	// data (hash)
	Data []byte
}

// Create Merkle tree
// txHashes: list of transaction hashes in the block
// The number of nodes at other levels other than the Merkle root node must be an even number, if it is an odd number, copy the last node
func NewMerkleTree(txHashes [][]byte) *MerkleTree {
	// Node list
	var nodes []MerkleNode
	// Determine the number of transaction data, if it is odd, copy the last one
	if len(txHashes)%2 != 0 {
		txHashes = append(txHashes, txHashes[len(txHashes)-1])
	}

	// Traverse all transaction data and generate leaf nodes by hashing
	for _, data := range txHashes {
		node := MakeMerkleNode(nil, nil, data)
		nodes = append(nodes, *node)
	}

	// create parent node from leaf node
	/*
	   Suppose there are 6 transactions, len(txHashes)=6
	   i = 0, len(nodes) = 4
	   i = 1, len(nodes) = 2
	   i = 2, len(nodes) = 1
	*/
	for i := 0; i < len(txHashes)/2; i++ {
		var parentNodes []MerkleNode // parent node list
		for j := 0; j < len(nodes); j += 2 {
			node := MakeMerkleNode(&nodes[j], &nodes[j+1], nil)
			parentNodes = append(parentNodes, *node)
		}

		if len(parentNodes)%2 != 0 {
			parentNodes = append(parentNodes, parentNodes[len(parentNodes)-1])
		}
		// In the end, only the hash value of the root node is stored in nodes
		nodes = parentNodes
	}
	mtree := MerkleTree{&nodes[0]}
	return &mtree
}

// Create Merkle node
func MakeMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := &MerkleNode{}
	// Determine leaf nodes
	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Data = hash[:]
	} else {
		// non-leaf node
		prveHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prveHashes)
		node.Data = hash[:]
	}
	// assignment of child nodes
	node.Left = left
	node.Right = right
	return node
}

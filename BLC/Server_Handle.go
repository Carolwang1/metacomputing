package BLC

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/labstack/gommon/log"
	"sort"
)

// request processing file management

//verion
func handleVersion(request []byte, bc *BlockChain) {
	fmt.Println("the request of version handle...")
	var buff bytes.Buffer
	var data Version
	// 1. Parse the request
	dataBytes := request[12:]
	// 2. Generate version structure
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)

	if err := decoder.Decode(&data); nil != err {
		log.Panicf("decode the version struct failed! %v\n", err)
	}
	// 3. Get the block height of the requester
	versionHeigth := data.Height
	// 4. Get the block height of its own node
	height := bc.GetHeigth()
	// If the block height of the current node is greater than versionHeigth
	// Send the current node version information to the requesting node
	fmt.Printf("client : %v, height : %v, versionHeigth : %v\n", data.AddrFrom, height, versionHeigth)
	if height > int64(versionHeigth) {
		sendVersion(data.AddrFrom, bc)
	} else if height < int64(versionHeigth) {
		// If the current node block height is less than versionHeigth
		// initiate a request to synchronize data to the sender
		sendGetBlocks(data.AddrFrom)
	}
}

// GetBlocks
// data synchronization request processing
func handleGetBlocks(request []byte, bc *BlockChain) {
	fmt.Println("the request of get blocks handle...")
	var buff bytes.Buffer
	var data GetBlocks
	// 1. Parse the request
	dataBytes := request[12:]
	// 2. Generate the getblocks structure
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)
	if err := decoder.Decode(&data); nil != err {
		log.Panicf("decode the getblocks struct failed! %v\n", err)
	}
	// 3. Get all block hashes of the block
	hashes := bc.GetBlockHases()
	sendInv(data.AddrFrom, hashes)
}

// Inv
func handleInv(request []byte, bc *BlockChain) {
	fmt.Println("the request of inv handle...")
	var buff bytes.Buffer
	var data Inv
	// 1. Parse the request
	dataBytes := request[12:]
	// 2. Generate Inv structure
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)
	if err := decoder.Decode(&data); nil != err {
		log.Panicf("decode the inv struct failed! %v\n", err)
	}
	var blocks [][]byte
	for _, hash := range data.Hashes {
		block := bc.GetBlock(hash)
		if block == nil {
			blocks = append(blocks, hash)
		}
	}
	sendGetData(data.AddrFrom, blocks)
	//for _, hash := range data.Hashes {
	// sendGetData(data.AddrFrom, hash)
	//}

	//for i := len(data.Hashes)-1; i >= 0; i-- {
	// sendGetData(data.AddrFrom,data.Hashes[i])
	//}

}

// GetData
// Process the request to get the specified block
func handleGetData(request []byte, bc *BlockChain) {
	fmt.Println("the request of get block handle...")
	var buff bytes.Buffer
	var data GetData
	// 1. Parse the request
	dataBytes := request[12:]
	// 2. Generate getData structure
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)
	if err := decoder.Decode(&data); nil != err {
		log.Panicf("decode the getData struct failed! %v\n", err)
	}
	var blockBytes [][]byte
	for _, hash := range data.ID {
		// 3. Obtain the block of the local node through the passed block hash
		blockByte := bc.GetBlock(hash)
		blockBytes = append(blockBytes, blockByte)
	}
	sendBlock(data.AddrFrom, blockBytes)

}

// Block
// When a new block is received, process it
func handleBlock(request []byte, bc *BlockChain) {
	fmt.Println("the request of handle block handle...")
	var buff bytes.Buffer
	var data BlockData
	// 1. Parse the request
	dataBytes := request[12:]
	// 2. Generate getData structure
	buff.Write(dataBytes)

	decoder := gob.NewDecoder(&buff)
	if err := decoder.Decode(&data); nil != err {
		log.Panicf("decode the blockdata struct failed! %v\n", err)
	}
	// 3. Add the received block to the blockchain
	blockBytes := data.Block
	Data := []*Block{}
	for _, blockByte := range blockBytes {
		block := DeserializeBlock(blockByte)
		Data = append(Data, block)
	}
	sort.Sort(Blocks(Data))

	for _, block := range Data {
		bc.AddBlock(block)
		// 4. Update utxo table
		utxoSet := UTXOSet{bc}
		utxoSet.update(block)
	}

}

type Blocks []*Block

func (s Blocks) Len() int {
	return len(s)
}
func (s Blocks) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Blocks) Less(i, j int) bool {
	return s[i].Heigth < s[j].Heigth
}

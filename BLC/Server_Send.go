package BLC

import (
	"bytes"
	"io"
	"log"
	"net"
)

// Request to send file

// send request
func sendMessage(to string, msg []byte) {
	// 1. connect to the server
	conn, err := net.Dial(PROTOCOL, to)
	if nil != err {
		log.Printf("connect to server [%s] failed! %v\n", to, err)
		return
	}
	if nil != conn {
		defer conn.Close()
	}
	// data to send
	_, err = io.Copy(conn, bytes.NewReader(msg))
	if nil != err {
		log.Printf("add the data to conn failed! %v\n", err)
		return
	}

}

// Blockchain version verification
func sendVersion(toAddress string, bc *BlockChain) {
	// 1. Get the block height of the current node
	height := bc.GetHeigth()
	//bc.DB.Close()
	// 2. Assemble the generated version
	versionData := Version{Height: int(height), AddrFrom: nodeAddress}
	// 3. data serialization
	data := gobEncode(versionData)
	// 4. Assemble commands and versions into complete requests
	request := append(commandToBytes(CMD_VERSION), data...)
	// 5. send request
	sendMessage(toAddress, request)
}

// Sync data from specified node
func sendGetBlocks(toAddress string) {
	// 1. Generate data
	data := gobEncode(GetBlocks{AddrFrom: nodeAddress})
	// 2. Assembly request
	request := append(commandToBytes(CMD_GETBLOCKS), data...)
	// 3. send request
	sendMessage(toAddress, request)
}

// Send a request to get the specified block
func sendGetData(toAddress string, hash [][]byte) {
	// 1. Generate data
	data := gobEncode(GetData{AddrFrom: nodeAddress, ID: hash})
	// 2. Assemble the request
	request := append(commandToBytes(CMD_GETDATA), data...)
	// 3. Send the request
	sendMessage(toAddress, request)
}

// show other nodes
func sendInv(toAddress string, hashes [][]byte) {
	// 1. Generate data
	data := gobEncode(Inv{AddrFrom: nodeAddress, Hashes: hashes})
	// 2. Assemble the request
	request := append(commandToBytes(CMD_INV), data...)
	// 3. Send the request
	sendMessage(toAddress, request)
}

// send block information
func sendBlock(toAddress string, block [][]byte) {
	// 1. Generate data
	data := gobEncode(BlockData{AddrFrom: nodeAddress, Block: block})
	// 2. Assemble the request
	request := append(commandToBytes(CMD_BLOCK), data...)
	// 3. Send the request
	sendMessage(toAddress, request)
}

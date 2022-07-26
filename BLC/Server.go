package BLC

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"net"
)

// Network service file management

// 3000 as the address of the boot node (master node)
var knownNodes = []string{"mcc.bootnode"}

// node address
var nodeAddress string

// start the service
func startServer(nodeID string) {
	fmt.Printf("Starting service[%s]...\n", nodeID)
	// Node address assignment
	nodeAddress = fmt.Sprintf("%s:8700", nodeID)
	// 1. Listen node
	listen, err := net.Listen(PROTOCOL, "0.0.0.0:8700")
	if nil != err {
		log.Panicf("listen address of %s failed! %v\n", nodeAddress, err)
	}

	defer listen.Close()
	// Get the blockchain object
	if nodeAddress != knownNodes[0] {
		bc := BlockchainObject(nodeID)
		sendVersion(knownNodes[0], bc)
		fmt.Println("Send request, sync data END...")
		bc.DB.Close()

	}

	for {
		// 2. Generate a connection and receive a request
		conn, err := listen.Accept()
		if nil != err {
			log.Printf("accept connect failed! %v\n", err)
		}
		// handle the request
		// Start a separate goroutine for request processing
		go handleConnection(conn, nodeID)
	}
}

// request handler
func handleConnection(conn net.Conn, nodeID string) {
	bc := BlockchainObject(nodeID)
	defer bc.DB.Close()
	request, err := ioutil.ReadAll(conn)
	if nil != err {
		log.Printf("Receive a Request failed! %v\n", err)
	}
	cmd := bytesToCommand(request[:12])
	fmt.Printf("Receive a Command: %s\n", cmd)
	switch cmd {
	case CMD_VERSION:
		handleVersion(request, bc)
	case CMD_GETDATA:
		handleGetData(request, bc)
	case CMD_GETBLOCKS:
		handleGetBlocks(request, bc)
	case CMD_INV:
		handleInv(request, bc)
	case CMD_BLOCK:
		handleBlock(request, bc)
	default:
		fmt.Println("Unknown command")
	}

}

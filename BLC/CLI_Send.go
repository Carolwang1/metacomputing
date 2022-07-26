package BLC

import (
	"fmt"
	"time"
)

// initiate transaction
func (cli *CLI) send(from, to, amount []string, nodeID string) {
	if !dbExist(nodeID) {
		fmt.Println("database does not exist...")
		return
	}
	// Get blockchain object
	blockchain := BlockchainObject(nodeID)
	nodeAddress = fmt.Sprintf("%s:8700", nodeID)
	if nodeAddress != knownNodes[0] {
		sendVersion(knownNodes[0], blockchain)
		time.Sleep(time.Duration(3) * time.Second)
	}
	defer blockchain.DB.Close()
	if len(from) != len(to) || len(from) != len(amount) {
		fmt.Println("The transaction parameters are entered incorrectly, please check the consistency...")
		return
	}
	// Initiate a transaction to generate a new block
	blockchain.MineNewBlock(from, to, amount, nodeID)
	//Call the function of utxo table to update utxo table
	utxoSet := &UTXOSet{Blockchain: blockchain}
	utxoSet.update(blockchain.Iterator().Next())

	if nodeAddress != knownNodes[0] {
		sendVersion(knownNodes[0], blockchain)
	}

}

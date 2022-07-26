package BLC

import "fmt"

// Check balances
func (cli *CLI) GetBalance(from string, nodeID string) {
	//Find the address UTXO
	//Get blockchain object
	blockchain := BlockchainObject(nodeID)
	defer blockchain.DB.Close() // close instance object
	utxoSet := UTXOSet{Blockchain: blockchain}
	amount := utxoSet.GetBalance(from)
	fmt.Printf("\taddress [%s] balance: [%d]\n", from, amount)

}

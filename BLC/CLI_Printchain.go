package BLC

import (
	"fmt"
	"os"
)

// Print complete blockchain information
func (cli *CLI) printchain(nodeID string) {
	if !dbExist(nodeID) {
		fmt.Println("database does not exist...")
		os.Exit(1)
	}

	blockchain := BlockchainObject(nodeID)
	fmt.Println("BlockchainObject end...")

	blockchain.PrintChain()
}

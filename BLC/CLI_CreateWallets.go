package BLC

import "fmt"

// Create a wallet collection from the command line
func (cli *CLI) CreateWallets(nodeID string) {
	// Create a wallet collection object
	wallets := NewWallets(nodeID)
	wallets.CreateWallet(nodeID)
	fmt.Printf("wallets : %v\n", wallets)
}

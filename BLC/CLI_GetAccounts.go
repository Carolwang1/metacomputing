package BLC

import "fmt"

// Get address list
func (cli *CLI) GetAccounts(nodeID string) {
	wallets := NewWallets(nodeID)
	fmt.Println("Account list")
	for key, _ := range wallets.Wallets {
		fmt.Printf("\t[%s]\n", key)
	}
}

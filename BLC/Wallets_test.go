package BLC

import (
	"fmt"
	"testing"
)

func TestWallets_CreateWallet(t *testing.T) {
	wallets := NewWallets("3001")
	wallets.CreateWallet("3001")
	fmt.Printf("wallets : %v\n", wallets.Wallets)
}

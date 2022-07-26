package BLC

import (
	"fmt"
	"testing"
)

func TestNewWallt(t *testing.T) {
	wallet := NewWallt()
	fmt.Printf("private key : %v\n", wallet.PrivateKey)
	fmt.Printf("public key : %v\n", wallet.PublicKey)
	fmt.Printf("wallet : %v\n", wallet)
}

func TestWallet_GetAddress(t *testing.T) {
	wallet := NewWallt()
	address := wallet.GetAddress()
	fmt.Printf("the address of coin is [%s]\n", address)
	fmt.Printf("the validation of current address is %v",
		IsValidForAddress([]byte(address)))
}

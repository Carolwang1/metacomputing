package BLC

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"os"
)

// The files managed by the wallet collection

// wallet collection persistence file
const walletFile = "Wallets_%s.dat"

// Implement the basic structure of the wallet collection
type Wallets struct {
	// key : address
	// value : wallet structure
	Wallets map[string]*Wallet
}

// Initialize wallet collection
func NewWallets(nodeID string) *Wallets {
	// Get wallet information from wallet file
	walletFile := fmt.Sprintf(walletFile, nodeID)
	fmt.Println(walletFile)
	// 1. First determine if the file exists
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		wallets := &Wallets{}
		wallets.Wallets = make(map[string]*Wallet)
		return wallets
	}
	// 2. The file exists, read the content
	fileContent, err := ioutil.ReadFile(walletFile)
	if nil != err {
		log.Panicf("read the file content failed! %v\n", err)
	}

	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if nil != err {
		log.Panicf("decode the file content failed! %v\n", err)
	}
	return &wallets
}

// add a new wallet to the collection
func (wallets *Wallets) CreateWallet(nodeID string) {
	// 1. Create wallet
	wallet := NewWallt()
	// 2. add
	wallets.Wallets[string(wallet.GetAddress())] = wallet
	// 3. Persist wallet information
	wallets.SaveWallets(nodeID)
}

// Persist wallet information (stored in a file)
func (w *Wallets) SaveWallets(nodeID string) {
	var content bytes.Buffer // wallet content
	// Register 256 ellipses. After registration, you can directly encode the curve interface internally
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(&w)

	if nil != err {
		log.Panicf("encode the struct of wallets failed %v\n", err)
	}
	walletFile := fmt.Sprintf(walletFile, nodeID)

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if nil != err {
		log.Panicf("write the content of wallet into file [%s] failed! %v\n", walletFile, err)
	}
}

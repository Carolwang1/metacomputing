package BLC

// Initialize the blockchain
func (cli *CLI) createBlockchain(address string, nodeID string) {
	bc := CreateBlockChainWithGenesisBlock(address, nodeID)
	defer bc.DB.Close()

	// set utxo reset action
	utxoSet := &UTXOSet{bc}
	utxoSet.ResetUTXOSet()
}

package BLC

func (cli *CLI) TestResetUTXO(nodeID string) {
	blockchain := BlockchainObject(nodeID)
	defer blockchain.DB.Close()
	utxoSet := UTXOSet{Blockchain: blockchain}
	utxoSet.ResetUTXOSet()
}

func (cli *CLI) TestFindUTXOMap() {

}

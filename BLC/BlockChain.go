package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/labstack/gommon/log"
	"math/big"
	"os"
	"strconv"
)

// Blockchain Management Documents
// Name database
const dbName = "block_%s.db"

// table name
const blockTableName = "blocks"

// Blockchain basic structure
type BlockChain struct {
	//Blocks []*Block		// slice of block
	DB *bolt.DB // database object

	Tip []byte // Save the hash value of the latest block
}

// Check if the database file exists
func dbExist(nodeID string) bool {
	// Generate database files for different nodes
	dbName := fmt.Sprintf(dbName, nodeID)
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		// 数据库文件不存在
		return false
	}
	return true
}

// Initialize the blockchain
func CreateBlockChainWithGenesisBlock(address string, nodeID string) *BlockChain {
	if dbExist(nodeID) {
		// The file already exists, indicating that the genesis block already exists
		fmt.Println("The genesis block already exists...")
		os.Exit(1)
	}
	// Save the hash value of the latest block
	var blockHash []byte
	// 1. Create or open a database
	dbName := fmt.Sprintf(dbName, nodeID)
	db, err := bolt.Open(dbName, 0600, nil)
	if nil != err {
		log.Panicf("create db [%s] failed %v\n", dbName, err)
	}
	// 2. Create a bucket and store the generated genesis block in the database
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b == nil {
			// no bucket found
			b, err := tx.CreateBucket([]byte(blockTableName))
			if nil != err {
				log.Panicf("create bucket [%s] failed %v\n", blockTableName, err)
			}
			// Generate a coinbase transaction
			txCoinbase := NewCoinbaseTransaction(address)
			// Generate genesis block
			genesisBlock := CreateGenesisBlock([]*Transaction{txCoinbase})
			// storage
			// 1. key,value What data represent--hash
			// 2. How to store the block structure in the database -- serialization
			err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if nil != err {
				log.Panicf("insert the genesis block failed %v\n", err)
			}
			blockHash = genesisBlock.Hash
			// Store the hash of the latest block
			// l : latest
			err = b.Put([]byte("l"), genesisBlock.Hash)
			if nil != err {
				log.Panicf("save the hash of genesis block failed %v\n", err)
			}
		}
		return nil
	})
	return &BlockChain{DB: db, Tip: blockHash}
}

// Traverse the database and output all block information
func (bc *BlockChain) PrintChain() {
	fmt.Println("Blockchain complete information...")
	var curBlock *Block
	bcit := bc.Iterator()
	// 循环读取
	for {
		fmt.Println("---------------------------------")
		curBlock = bcit.Next()
		fmt.Printf("\tHash:%x\n", curBlock.Hash)
		fmt.Printf("\tPrevBlockHash:%x\n", curBlock.PrevBlockHash)
		fmt.Printf("\tTimeStamp:%v\n", curBlock.TimeStamp)
		fmt.Printf("\tHeigth:%d\n", curBlock.Heigth)
		fmt.Printf("\tNonce:%d\n", curBlock.Nonce)
		fmt.Printf("\tTxs:%v\n", curBlock.Txs)
		for _, tx := range curBlock.Txs {
			fmt.Printf("\t\ttx-hash : %x\n", tx.TxHash)
			fmt.Printf("\t\tinput...\n")
			for _, vin := range tx.Vins {
				fmt.Printf("\t\t\tvin-txHash : %x\n", vin.TxHash)
				fmt.Printf("\t\t\tvin-vout : %v\n", vin.Vout)
				fmt.Printf("\t\t\tvin-PublicKey : %x\n", vin.PublicKey)
				fmt.Printf("\t\t\tvin-Signature : %x\n", vin.Signature)
			}
			fmt.Printf("\t\toutput...\n")
			for _, vout := range tx.Vouts {
				fmt.Printf("\t\t\tvout-value:%d\n", vout.Value)
				fmt.Printf("\t\t\tvout-Ripemd160Hash:%x\n", vout.Ripemd160Hash)
			}
		}
		// Exit conditions
		// convert to big.int
		var hashInt big.Int
		hashInt.SetBytes(curBlock.PrevBlockHash)
		// compare
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			// Traverse to the genesis block
			break
		}
	}
}

// Get a blockchain object
func BlockchainObject(nodeID string) *BlockChain {
	// Get DB
	dbName := fmt.Sprintf(dbName, nodeID)
	db, err := bolt.Open(dbName, 0600, nil)
	if nil != err {
		log.Panicf("open the db [%s] failed! %v\n", dbName, err)
	}

	// Get Tips
	var tip []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if nil != b {
			tip = b.Get([]byte("l"))
		}
		return nil
	})
	if nil != err {
		log.Panicf("get the blockchain object failed! %v\n", err)
	}

	return &BlockChain{DB: db, Tip: tip}
}

// Implement mining function
// By receiving transactions, blocks are generated
func (blockchain *BlockChain) MineNewBlock(from, to, amount []string, nodeID string) {
	// Shelving Transaction Generation Steps
	var block *Block
	var txs []*Transaction
	// Participants in the traversal transaction
	for index, address := range from {
		value, _ := strconv.Atoi(amount[index])
		if value <= 0 {
			fmt.Println("The number of transactions must be greater than 0")
			os.Exit(1)
		}
		// generate new transaction
		tx := NewSimpleTransaction(address, to[index], value, blockchain, txs, nodeID)
		// sign
		// Append to the transaction list of txs
		txs = append(txs, tx)
		// Give a certain reward to the initiator (miner) of the transaction
		tx = NewCoinbaseTransaction(address)
		txs = append(txs, tx)
	}

	// Get the latest block from the database
	blockchain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if nil != b {
			// Get the latest block hash
			hash := b.Get([]byte("l"))
			// Get the latest block
			blockBytes := b.Get(hash)
			// deserialize
			block = DeserializeBlock(blockBytes)
		}
		return nil
	})
	// Verify transaction signature here
	// Verify the signature of every transaction in txs
	for _, tx := range txs {
		// Verify the signature, as long as the verification of a signature fails, panic
		if blockchain.VerifyTransaction(tx) == false {
			log.Panicf("ERROR : tx [%x] verify failed!\n")
		}
	}
	// Generate a new block (packaging) from the latest block in the database
	block = NewBlock(block.Heigth+1, block.Hash, txs)
	// Persist newly generated blocks to the database
	blockchain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if nil != b {
			err := b.Put(block.Hash, block.Serialize())
			if nil != err {
				log.Panicf("update the new block to db failed %v\n", err)
			}

			// Update the hash of the latest block
			err = b.Put([]byte("l"), block.Hash)
			if nil != err {
				log.Panicf("update the latest block hash to db failed %v\n", err)
			}
			blockchain.Tip = block.Hash
		}
		return nil
	})
}

// Get all the spent output of the specified address
func (blockchain *BlockChain) SpentOutputs(address string) map[string][]int {
	// Spent output cache
	spentTXOutputs := make(map[string][]int)
	// get iterator object
	bcit := blockchain.Iterator()
	for {
		block := bcit.Next()
		for _, tx := range block.Txs {
			// Exclude coinbase transactions
			if !tx.IsCoinbaseTransaction() {
				for _, in := range tx.Vins {
					if in.UnLockRipemd160Hash(StringToHash160(address)) {
						key := hex.EncodeToString(in.TxHash)
						// add to the cache of spent output
						spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
					}

				}
			}
		}

		// exit loop condition
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	return spentTXOutputs
}

// Find the UTXO of the specified address
/*
	Traverse to find every transaction in every block in the blockchain database
	Find every output in every transaction
	Determine whether each output satisfies the following conditions
	1. Belong to the incoming address
	2. Is it not spent
	1. First, traverse the blockchain database once and store all the spent OUTPUT in a cache
	2. Traverse the blockchain database again and check if each VOUT is contained in the previous cache of spent outputs
*/
func (blockchain *BlockChain) UnUTXOS(address string, txs []*Transaction) []*UTXO {
	// 1. Traverse the database to find all transactions related to address
	// get iterator
	bcit := blockchain.Iterator()
	// List of unspent output for current address
	var unUTXOS []*UTXO
	// Get all the spent output of the specified address
	spentTXOutputs := blockchain.SpentOutputs(address)
	// cache iteration
	// Find spent output in cache
	for _, tx := range txs {
		// Judge coinbaseTransaction
		if !tx.IsCoinbaseTransaction() {
			for _, in := range tx.Vins {

				// Judge the user
				if in.UnLockRipemd160Hash(StringToHash160(address)) {
					// add to the map of spent output
					key := hex.EncodeToString(in.TxHash)
					spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
				}

			}
		}
	}
	// Traverse the UTXOs in the cache
	for _, tx := range txs {
		// Add a jump to cache output
	WorkCacheTx:
		for index, vout := range tx.Vouts {
			if vout.UnLockScriptPubkeyWithAddress(address) {
				//if vout.CheckPubkeyWithAddress(address) {
				if len(spentTXOutputs) != 0 {
					var isUtxoTx bool // Determine whether a transaction is referenced by other transactions
					for txHash, indexArray := range spentTXOutputs {
						txHashStr := hex.EncodeToString(tx.TxHash)
						if txHash == txHashStr {
							// The currently traversed transaction already has outputs referenced by the inputs of other transactions
							isUtxoTx = true
							// Add a state variable to determine whether the specified output is referenced
							var isSpentUTXO bool
							for _, voutIndex := range indexArray {
								if index == voutIndex {
									// This output is quoted
									isSpentUTXO = true
									// Jump out of the current vout judgment logic and proceed to the next output judgment
									continue WorkCacheTx
								}
							}
							if isSpentUTXO == false {
								utxo := &UTXO{tx.TxHash, index, vout}
								unUTXOS = append(unUTXOS, utxo)
							}
						}
					}
					if isUtxoTx == false {
						// Indicates that all address-related outputs in the current transaction are UTXO
						utxo := &UTXO{tx.TxHash, index, vout}
						unUTXOS = append(unUTXOS, utxo)
					}
				} else {
					utxo := &UTXO{tx.TxHash, index, vout}
					unUTXOS = append(unUTXOS, utxo)
				}
			}
		}
	}

	// First traverse the UTXO in the cache, if the balance is enough, return directly, if not, then traverse the UTXO in the db file
	// Iterate over the database, keep getting the next block
	for {
		block := bcit.Next()
		// Iterate over each transaction in the block
		for _, tx := range block.Txs {
			// jump
		work:
			for index, vout := range tx.Vouts {
				// index:The current output is at the mid-index position of the current transaction
				// vout:current output
				if vout.UnLockScriptPubkeyWithAddress(address) {
					//if vout.CheckPubkeyWithAddress(address) {
					// The current vout belongs to the incoming address
					if len(spentTXOutputs) != 0 {
						var isSpentOutput bool // default false
						for txHash, indexArray := range spentTXOutputs {
							for _, i := range indexArray {
								// txHash : The transaction hash referenced by the current output
								// indexArray: A list of vout indices associated with the hash
								if txHash == hex.EncodeToString(tx.TxHash) && index == i {
									// txHash == hex.EncodeToString(tx.TxHash),
									// Indicates that the current transaction tx has at least an output referenced by the input of other transactions
									// index == i It means that the current output is referenced by other transactions
									// Jump to the outermost loop to judge the next VOUT
									isSpentOutput = true
									continue work
								}
							}
						}
						/*
							// UTXO structure management
								type UTXO struct {
									// UTXO corresponding transaction hash
									TxHash 		[]byte
									// The index of the UTXO in the output list of the transaction to which it belongs
									Index		int
									// Output itself
									Output 		*TxOutput
								}
						*/
						if isSpentOutput == false {
							utxo := &UTXO{tx.TxHash, index, vout}
							unUTXOS = append(unUTXOS, utxo)
						}
					} else {
						// Add all outputs from the current address to the unspent outputs
						utxo := &UTXO{tx.TxHash, index, vout}
						unUTXOS = append(unUTXOS, utxo)
					}
				}
			}
		}

		// exit loop condition
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	return unUTXOS
}

// Check balances
func (blockchain *BlockChain) getBalance(address string) int {
	var amount int
	utxos := blockchain.UnUTXOS(address, []*Transaction{})
	for _, utxo := range utxos {
		amount += utxo.Output.Value
	}
	return amount
}

// Find the available UTXO at the specified address, and interrupt the search if the amount is exceeded
// Update the number of UTXOs at the specified address in the current database
// txs: The list of transactions in the cache (for multi-transaction processing)
func (blockchain *BlockChain) FindSpendableUTXO(from string,
	amount int, txs []*Transaction) (int, map[string][]int) {
	// Available UTXOs
	spendableUTXO := make(map[string][]int)

	var value int
	utxos := blockchain.UnUTXOS(from, txs)
	// Traverse UTXOs
	for _, utxo := range utxos {
		if utxo.Output.Value <= 0 {
			continue
		}
		value += utxo.Output.Value
		// Calculate transaction hash
		hash := hex.EncodeToString(utxo.TxHash)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)
		if value >= amount {
			break
		}
	}

	// all traversal is completed, still less than amount
	// insufficient funds
	if value < amount {
		fmt.Printf("address [%s] Insufficient balance, current balance [%d], transfer amount [%d]\n", from, value, amount)
		os.Exit(1)
	}

	return value, spendableUTXO
}

// Find a transaction by the specified transaction hash
func (blockchain *BlockChain) FindTransaction(ID []byte) Transaction {
	bcit := blockchain.Iterator()
	for {
		block := bcit.Next()
		for _, tx := range block.Txs {
			if bytes.Compare(ID, tx.TxHash) == 0 {
				// 找到该交易
				return *tx
			}
		}

		// quit
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}

	fmt.Printf("No deal found[%x]\n", ID)
	return Transaction{}
}

// transaction signature
func (blockchain *BlockChain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	// coinbase transactions do not require signatures
	if tx.IsCoinbaseTransaction() {
		return
	}
	// Process the input of the transaction and find the transaction to which the vout referenced by the input in the tx belongs (find the sender)
	// Sign every UTXO we spend
	// store the referenced transaction
	prevTxs := make(map[string]Transaction)
	for _, vin := range tx.Vins {
		// Find the transaction referenced by the current transaction input
		tx := blockchain.FindTransaction(vin.TxHash)
		prevTxs[hex.EncodeToString(tx.TxHash)] = tx
	}
	// sign
	tx.Sign(privKey, prevTxs)
}

// Verify signature
func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbaseTransaction() {
		return true
	}
	prevTxs := make(map[string]Transaction)
	// 查找输入引用的交易
	for _, vin := range tx.Vins {
		tx := bc.FindTransaction(vin.TxHash)
		prevTxs[hex.EncodeToString(tx.TxHash)] = tx
	}

	return tx.Verify(prevTxs)
}

// Exit conditions
func isBreakLoop(prevBlockHash []byte) bool {
	var hashInt big.Int
	hashInt.SetBytes(prevBlockHash)
	if hashInt.Cmp(big.NewInt(0)) == 0 {
		return true
	}
	return false
}

// Find all spent outputs across the entire blockchain
func (blockchain *BlockChain) FindAllSpentOutputs() map[string][]*TxInput {
	bcit := blockchain.Iterator()
	// Store spent output
	spentTXOutputs := make(map[string][]*TxInput)
	for {
		block := bcit.Next()
		for _, tx := range block.Txs {
			if !tx.IsCoinbaseTransaction() {
				for _, txInput := range tx.Vins {
					txHash := hex.EncodeToString(txInput.TxHash)
					spentTXOutputs[txHash] = append(spentTXOutputs[txHash], txInput)
				}
			}
		}

		if isBreakLoop(block.PrevBlockHash) {
			break
		}
	}
	return spentTXOutputs
}

// Find UTXOs for all addresses in the entire blockchain
func (blockchain *BlockChain) FindUTXOMap() map[string]*TXOutputs {
	// Traverse the blockchain
	bcit := blockchain.Iterator()
	// output collection
	utxoMaps := make(map[string]*TXOutputs)
	// Find spent output
	spentTXOutputs := blockchain.FindAllSpentOutputs()

	for {
		block := bcit.Next()
		for _, tx := range block.Txs {
			txOutputs := &TXOutputs{[]*TxOutput{}}
			txHash := hex.EncodeToString(tx.TxHash)
			// Get vouts for each transaction
		WorkOutLoop:
			for index, vout := range tx.Vouts {
				// Get the input of the specified transaction
				txInputs := spentTXOutputs[txHash]
				if len(txInputs) > 0 {
					isSpent := false
					for _, in := range txInputs {
						// Find the owner of the specified output
						outPubkey := vout.Ripemd160Hash
						inPubkey := in.PublicKey
						if bytes.Compare(outPubkey, Ripemd160Hash(inPubkey)) == 0 {
							if index == in.Vout {
								isSpent = true
								continue WorkOutLoop
							}
						}
					}

					if isSpent == false {
						// The current output is not included in txInputs
						txOutputs.TXOutputs = append(txOutputs.TXOutputs, vout)
					}
				} else {
					// If no input refers to the output of the transaction, it means that all the outputs in the current transaction are UTXO
					txOutputs.TXOutputs = append(txOutputs.TXOutputs, vout)
				}
			}
			utxoMaps[txHash] = txOutputs
		}

		if isBreakLoop(block.PrevBlockHash) {
			break
		}
	}
	return utxoMaps
}

// Get the block height of the current block
func (bc *BlockChain) GetHeigth() int64 {
	return bc.Iterator().Next().Heigth
}

// Get all block hashes of the blockchain
func (bc *BlockChain) GetBlockHases() [][]byte {
	var blockHashes [][]byte
	bcit := bc.Iterator()
	for {
		block := bcit.Next()
		blockHashes = append(blockHashes, block.Hash)
		if isBreakLoop(block.PrevBlockHash) {
			break
		}
	}
	return blockHashes
}

// Get the block data of the specified hash
func (bc *BlockChain) GetBlock(hash []byte) []byte {
	var blockByte []byte
	bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if nil != b {
			blockByte = b.Get(hash)
		}
		return nil
	})
	return blockByte
}

// add block
func (bc *BlockChain) AddBlock(block *Block) {
	err := bc.DB.Update(func(tx *bolt.Tx) error {
		// 1. get datasheet
		b := tx.Bucket([]byte(blockTableName))
		if nil != b {
			// Determine whether the block that needs to be passed in already exists
			if b.Get(block.Hash) != nil {
				// 已经存在，不需要添加
				return nil
			}
			// does not exist, add to database
			err := b.Put(block.Hash, block.Serialize())
			if nil != err {
				log.Panicf("sync the block failed! %v\n", err)
			}
			blockHash := b.Get([]byte("l"))
			latesBlock := b.Get(blockHash)
			rawBlock := DeserializeBlock(latesBlock)
			if rawBlock.Heigth < block.Heigth {
				b.Put([]byte("l"), block.Hash)
				bc.Tip = block.Hash
			}
		}
		return nil
	})
	if nil != err {
		log.Panicf("update the db when insert the new block failed! %v\n", err)
	}
	fmt.Println("the new block is added!")
}

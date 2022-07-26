package BLC

import (
	"flag"
	"fmt"
	"github.com/labstack/gommon/log"
	"os"
)

// Manage the command line operations of the blockchain

// client object
type CLI struct{}

// usage show
func PrintUsage() {
	fmt.Println("Usage:")
	// Initialize the blockchain
	fmt.Printf("\tcreateblockchain -address address -- Create a blockchain\n")
	// add block
	//fmt.Printf("\taddblock -data DATA -- add block\n")
	// Print complete block information
	fmt.Printf("\tprintchain -- output blockchain information\n")
	// Transfer by command
	fmt.Printf("\tsend -from FROM -to TO -amount AMOUNT -- Initiate a transfer\n")
	fmt.Printf("\tDescription of transfer parameters\n")
	fmt.Printf("\t\t-from FROM -- Transfer source address\n")
	fmt.Printf("\t\t-to TO -- Transfer destination address\n")
	fmt.Printf("\t\t-AMOUNT amount -- transfer amount\n")
	// Check balances
	fmt.Printf("\tgetbalance -address FROM -- Query the balance of the specified address\n")
	fmt.Println("\tQuery balance parameter description")
	fmt.Printf("\t\t-address -- Address to check balance\n")
	// Wallet management
	fmt.Printf("\tcreatewallet -- Create wallet\n")
	fmt.Printf("\taccounts -- Get a list of wallet addresses\n")
	fmt.Printf("\tutxo -method METHOD -- Test the method specified in the UTXO Table function\n")
	fmt.Printf("\t\tMETHOD -- method name\n")
	fmt.Printf("\t\t\treset -- reset UTXOtable\n")
	fmt.Printf("\t\t\tbalance - Find all UTXOs\n")
	fmt.Printf("\tset_id -port PORT -- set node number\n")
	fmt.Printf("\t\tport -- Visited node number\n")
	fmt.Printf("\tstart -- start node service\n")
}

// add block
func (cli *CLI) addBlock(txs []*Transaction) {

}

// command line run function
func (cli *CLI) Run() {
	nodeId := GetEnvNodeId()
	// Number of detection parameters
	IsValidArgs()
	// New related command
	// add block
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	// Output the complete information of the blockchain
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	// Create a blockchain
	createBLCWithGenesisBlockCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	// initiate transaction
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	// Check balance command
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	// Wallet management related commands
	// Create a wallet collection
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	// Get address list
	getAccountsCmd := flag.NewFlagSet("accounts", flag.ExitOnError)
	// utxo test command
	UTXOTestCmd := flag.NewFlagSet("utxo", flag.ExitOnError)
	// Node number setting command
	setNodeIdCmd := flag.NewFlagSet("set_id", flag.ExitOnError)
	// Node service start command
	startNodeCmd := flag.NewFlagSet("start", flag.ExitOnError)
	// data parameter processing
	// add block parameters
	flagAddBlockArg := addBlockCmd.String("data", "sent 100 btc to player", "add block data")
	// Miner address (receive reward) parameter specified when creating the blockchain
	flagCreateBlockchainArg := createBLCWithGenesisBlockCmd.String("address",
		"troytan", "Specify the miner's address to receive system rewards")
	// Initiate transaction parameters
	flagSendFromArg := sendCmd.String("from", "", "Transfer source address")
	flagSendToArg := sendCmd.String("to", "", "Transfer destination address")
	flagSendAmountArg := sendCmd.String("amount", "", "transfer amount")
	// Query balance command line parameters
	flagGetBalanceArg := getBalanceCmd.String("address", "", "address to query")
	// UTXO test command line arguments
	flagUTXOArg := UTXOTestCmd.String("method", "", "UTXO Table Related operations\n")
	// port number parameter
	flagPortArg := setNodeIdCmd.String("port", "", "set node ID")
	// Judgment command
	switch os.Args[1] {
	case "start":
		if err := startNodeCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse cmd start node server failed! %v\n", err)
		}
	case "set_id":
		if err := setNodeIdCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse cmd set node id failed! %v\n", err)
		}
	case "utxo":
		if err := UTXOTestCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse cmd operate utxo table failed! %v\n", err)
		}
	case "accounts":
		if err := getAccountsCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse cmd get accounts failed! %v\n", err)
		}
	case "createwallet":
		if err := createWalletCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse cmd create wallet failed! %v\n", err)
		}
	case "getbalance":
		if err := getBalanceCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse cmd get balance failed! %v\n", err)
		}
	case "send":
		if err := sendCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse sendCmd failed! %v\n", err)
		}
	case "addblock":
		if err := addBlockCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse addBlockCmd failed! %v\n", err)
		}
	case "printchain":
		if err := printChainCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse printchainCmd failed! %v\n", err)
		}
	case "createblockchain":
		if err := createBLCWithGenesisBlockCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse createBLCWithGenesisBlockCmd failed! %v\n", err)
		}
	default:
		// No command was passed or the command passed was not in the command list above
		PrintUsage()
		os.Exit(1)

	}

	// Node start service
	if startNodeCmd.Parsed() {
		cli.startNode(nodeId)
	}

	// Node ID settings
	if setNodeIdCmd.Parsed() {
		if *flagPortArg == "" {
			fmt.Println("Please enter the port number...")
			os.Exit(1)
		}
		cli.SetNodeId(*flagPortArg)
	}

	// utxo table operate
	if UTXOTestCmd.Parsed() {
		switch *flagUTXOArg {
		case "balance":
			cli.TestFindUTXOMap()
		case "reset":
			cli.TestResetUTXO(nodeId)
		default:
		}
	}

	// Get address list
	if getAccountsCmd.Parsed() {
		cli.GetAccounts(nodeId)
	}

	// Create wallet
	if createWalletCmd.Parsed() {
		cli.CreateWallets(nodeId)
	}

	// Check balances
	if getBalanceCmd.Parsed() {
		if *flagGetBalanceArg == "" {
			fmt.Println("Please enter the inquiry address...")
			os.Exit(1)
		}
		cli.GetBalance(*flagGetBalanceArg, nodeId)
	}
	// Initiate a transfer
	if sendCmd.Parsed() {
		if *flagSendFromArg == "" {
			fmt.Println("Source address cannot be empty...")
			PrintUsage()
			os.Exit(1)
		}
		if *flagSendToArg == "" {
			fmt.Println("Destination address cannot be empty...")
			PrintUsage()
			os.Exit(1)
		}
		if *flagSendAmountArg == "" {
			fmt.Println("The transfer amount cannot be empty...")
			PrintUsage()
			os.Exit(1)
		}
		fmt.Printf("\tFROM:[%s]\n", JSONToSlice(*flagSendFromArg))
		fmt.Printf("\tTO:[%s]\n", JSONToSlice(*flagSendToArg))
		fmt.Printf("\tAMOUNT:[%s]\n", JSONToSlice(*flagSendAmountArg))
		cli.send(JSONToSlice(*flagSendFromArg), JSONToSlice(*flagSendToArg), JSONToSlice(*flagSendAmountArg), nodeId)
	}
	// add block command
	if addBlockCmd.Parsed() {
		if *flagAddBlockArg == "" {
			PrintUsage()
			os.Exit(1)
		}
		cli.addBlock([]*Transaction{})
	}
	// output blockchain information
	fmt.Println("nodeID : ", nodeId)
	if printChainCmd.Parsed() {
		cli.printchain(nodeId)
	}
	// Create blockchain command
	if createBLCWithGenesisBlockCmd.Parsed() {
		if *flagCreateBlockchainArg == "" {
			PrintUsage()
			os.Exit(1)
		}
		cli.createBlockchain(*flagCreateBlockchainArg, nodeId)
	}
}

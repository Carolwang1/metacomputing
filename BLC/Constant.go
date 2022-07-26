package BLC

// web service constant management
// protocol
const PROTOCOL = "tcp"

// command length
const CMMAND_LENGTH = 12

// command classification
const (
	// Verify that the current node end block is the latest block
	CMD_VERSION = "version"
	// Get the block from the longest chain
	CMD_GETBLOCKS = "getblocks"
	// Show other nodes which blocks the current node has
	CMD_INV = "inv"
	// request the specified block
	CMD_GETDATA = "getdata"
	// After receiving the new block, process it
	CMD_BLOCK = "block"
)

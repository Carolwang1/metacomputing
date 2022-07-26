package BLC

// Current block version information (determines whether the block needs to be synchronized)
type Version struct {
	//Version int // version number
	Height   int    // The block height of the current node
	AddrFrom string // the address of the current node
}

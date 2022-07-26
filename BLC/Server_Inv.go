package BLC

// exhibit
type Inv struct {
	AddrFrom string   // the address of the current node
	Hashes   [][]byte // Hash list of blocks owned by the current node
}

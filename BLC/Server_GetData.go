package BLC

// Request a specified block
type GetData struct {
	AddrFrom string   // current address
	ID       [][]byte // block hash
}

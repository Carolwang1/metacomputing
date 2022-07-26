package BLC

type BlockData struct {
	AddrFrom string   // node address
	Block    [][]byte // block data (serialized data)
}

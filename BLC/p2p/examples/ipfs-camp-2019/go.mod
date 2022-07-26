module github.com/libp2p/go-libp2p/examples/ipfs-camp-2019

go 1.12

require (
	github.com/gogo/protobuf v1.3.2
	github.com/libp2p/go-libp2p v0.13.0
	github.com/libp2p/go-libp2p-core v0.8.5
	github.com/libp2p/go-libp2p-discovery v0.5.0
	github.com/libp2p/go-libp2p-kad-dht v0.11.1
	github.com/libp2p/go-libp2p-mplex v0.4.1
	github.com/libp2p/go-libp2p-pubsub v0.4.1
	github.com/libp2p/go-libp2p-secio v0.2.2
	github.com/libp2p/go-libp2p-yamux v0.5.4
	github.com/libp2p/go-tcp-transport v0.2.3
	github.com/libp2p/go-ws-transport v0.4.0
	github.com/multiformats/go-multiaddr v0.3.3
)

// Ensure that examples always use the go-libp2p version in the same git checkout.
replace github.com/libp2p/go-libp2p => ./../..

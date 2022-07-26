package backpressure_tests

import (
	"context"
	"os"
	"testing"
	"time"

	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	"github.com/stretchr/testify/require"

	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/network"
	protocol "github.com/libp2p/go-libp2p-core/protocol"
	swarmt "github.com/libp2p/go-libp2p-swarm/testing"
)

var log = logging.Logger("backpressure")

// TestStBackpressureStreamWrite tests whether streams see proper
// backpressure when writing data over the network streams.
func TestStBackpressureStreamWrite(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	h1 := bhost.New(swarmt.GenSwarm(t, ctx))
	h2 := bhost.New(swarmt.GenSwarm(t, ctx))

	// setup sender handler on 1
	h1.SetStreamHandler(protocol.TestingID, func(s network.Stream) {
		defer s.Reset()
		<-ctx.Done()
	})

	h2pi := h2.Peerstore().PeerInfo(h2.ID())
	log.Debugf("dialing %s", h2pi.Addrs)
	if err := h1.Connect(ctx, h2pi); err != nil {
		t.Fatal("Failed to connect:", err)
	}

	// open a stream, from 2->1, this is our reader
	s, err := h2.NewStream(ctx, h1.ID(), protocol.TestingID)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Reset()

	// If nobody is reading, we should eventually timeout.
	require.NoError(t, s.SetWriteDeadline(time.Now().Add(100*time.Millisecond)))
	data := make([]byte, 16*1024)
	for i := 0; i < 5*1024; i++ { // write at most 100MiB
		_, err := s.Write(data)
		if err != nil {
			require.True(t, os.IsTimeout(err), err)
			return
		}
	}
	t.Fatal("should have timed out")
}

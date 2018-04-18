package coremock

import (
	"context"
	"net"

	commands "github.com/ipsn/go-ipfs/commands"
	core "github.com/ipsn/go-ipfs/core"
	"github.com/ipsn/go-ipfs/repo"
	config "github.com/ipsn/go-ipfs/repo/config"

	mocknet "github.com/ipsn/go-ipfs/gxlibs/github.com/libp2p/go-libp2p/p2p/net/mock"
	host "github.com/ipsn/go-ipfs/gxlibs/github.com/libp2p/go-libp2p-host"
	testutil "github.com/ipsn/go-ipfs/gxlibs/github.com/libp2p/go-testutil"
	datastore "github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-datastore"
	syncds "github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-datastore/sync"
	pstore "github.com/ipsn/go-ipfs/gxlibs/github.com/libp2p/go-libp2p-peerstore"
	smux "github.com/ipsn/go-ipfs/gxlibs/github.com/libp2p/go-stream-muxer"
	ipnet "github.com/ipsn/go-ipfs/gxlibs/github.com/libp2p/go-libp2p-interface-pnet"
	peer "github.com/ipsn/go-ipfs/gxlibs/github.com/libp2p/go-libp2p-peer"
	metrics "github.com/ipsn/go-ipfs/gxlibs/github.com/libp2p/go-libp2p-metrics"
)

// NewMockNode constructs an IpfsNode for use in tests.
func NewMockNode() (*core.IpfsNode, error) {
	ctx := context.Background()

	// effectively offline, only peer in its network
	return core.NewNode(ctx, &core.BuildCfg{
		Online: true,
		Host:   MockHostOption(mocknet.New(ctx)),
	})
}

func MockHostOption(mn mocknet.Mocknet) core.HostOption {
	return func(ctx context.Context, id peer.ID, ps pstore.Peerstore, bwr metrics.Reporter, fs []*net.IPNet, _ smux.Transport, _ ipnet.Protector, _ *core.ConstructPeerHostOpts) (host.Host, error) {
		return mn.AddPeerWithPeerstore(id, ps)
	}
}

func MockCmdsCtx() (commands.Context, error) {
	// Generate Identity
	ident, err := testutil.RandIdentity()
	if err != nil {
		return commands.Context{}, err
	}
	p := ident.ID()

	conf := config.Config{
		Identity: config.Identity{
			PeerID: p.String(),
		},
	}

	r := &repo.Mock{
		D: syncds.MutexWrap(datastore.NewMapDatastore()),
		C: conf,
	}

	node, err := core.NewNode(context.Background(), &core.BuildCfg{
		Repo: r,
	})
	if err != nil {
		return commands.Context{}, err
	}

	return commands.Context{
		Online:     true,
		ConfigRoot: "/tmp/.mockipfsconfig",
		LoadConfig: func(path string) (*config.Config, error) {
			return &conf, nil
		},
		ConstructNode: func() (*core.IpfsNode, error) {
			return node, nil
		},
	}, nil
}
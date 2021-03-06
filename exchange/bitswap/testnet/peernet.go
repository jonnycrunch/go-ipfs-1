package bitswap

import (
	"context"

	bsnet "github.com/ipsn/go-ipfs/exchange/bitswap/network"

	mockpeernet "github.com/ipsn/go-ipfs/gxlibs/github.com/libp2p/go-libp2p/p2p/net/mock"
	testutil "github.com/ipsn/go-ipfs/gxlibs/github.com/libp2p/go-testutil"
	ds "github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-datastore"
	mockrouting "github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-ipfs-routing/mock"
	peer "github.com/ipsn/go-ipfs/gxlibs/github.com/libp2p/go-libp2p-peer"
)

type peernet struct {
	mockpeernet.Mocknet
	routingserver mockrouting.Server
}

func StreamNet(ctx context.Context, net mockpeernet.Mocknet, rs mockrouting.Server) (Network, error) {
	return &peernet{net, rs}, nil
}

func (pn *peernet) Adapter(p testutil.Identity) bsnet.BitSwapNetwork {
	client, err := pn.Mocknet.AddPeer(p.PrivateKey(), p.Address())
	if err != nil {
		panic(err.Error())
	}
	routing := pn.routingserver.ClientWithDatastore(context.TODO(), p, ds.NewMapDatastore())
	return bsnet.NewFromIpfsHost(client, routing)
}

func (pn *peernet) HasPeer(p peer.ID) bool {
	for _, member := range pn.Mocknet.Peers() {
		if p == member {
			return true
		}
	}
	return false
}

var _ Network = (*peernet)(nil)

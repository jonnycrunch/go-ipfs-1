package ipns

import (
	"context"

	"github.com/ipsn/go-ipfs/core"
	nsys "github.com/ipsn/go-ipfs/namesys"
	path "github.com/ipsn/go-ipfs/path"
	ft "github.com/ipsn/go-ipfs/unixfs"
	ci "github.com/ipsn/go-ipfs/gxlibs/github.com/libp2p/go-libp2p-crypto"
)

// InitializeKeyspace sets the ipns record for the given key to
// point to an empty directory.
func InitializeKeyspace(n *core.IpfsNode, key ci.PrivKey) error {
	ctx, cancel := context.WithCancel(n.Context())
	defer cancel()

	emptyDir := ft.EmptyDirNode()

	err := n.Pinning.Pin(ctx, emptyDir, false)
	if err != nil {
		return err
	}

	err = n.Pinning.Flush()
	if err != nil {
		return err
	}

	pub := nsys.NewRoutingPublisher(n.Routing, n.Repo.Datastore())

	return pub.Publish(ctx, key, path.FromCid(emptyDir.Cid()))
}

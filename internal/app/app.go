package app

import (
	"context"
	"github.com/akbariandev/jumpy/internal/chain"
	"github.com/akbariandev/jumpy/internal/p2p"
	_ "github.com/libp2p/go-libp2p/p2p/host/peerstore"
)

const defaultHostGroupName = "jumpy"

type Application struct {
	nodes []*p2p.PeerStream
}

func (a *Application) ListNodes() []*p2p.PeerStream {
	return a.nodes
}

func (a *Application) Start(nodesCount int, groupName string) {
	ctx := context.Background()
	if len(groupName) == 0 {
		groupName = defaultHostGroupName
	}

	port := 2000
	for nodesCount > 0 {
		ps, err := p2p.NewPeerStream(port+nodesCount, approveFunc)
		if err != nil {
			panic(err)
		}
		ps.Run(ctx, groupName, false)
		a.nodes = append(a.nodes, ps)
		nodesCount--
	}

	select {}
}

func approveFunc(block *chain.Block) bool {
	if block != nil {
		return true
	}
	return false
}

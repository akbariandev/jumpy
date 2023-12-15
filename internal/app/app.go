package app

import (
	"context"
	"github.com/akbariandev/jumpy/internal/p2p"
	_ "github.com/libp2p/go-libp2p/p2p/host/peerstore"
)

const defaultHostGroupName = "jumpy"

func Start(listenPort int, hostGroupName string) {
	ctx := context.Background()

	if len(hostGroupName) == 0 {
		hostGroupName = defaultHostGroupName
	}

	ps, err := p2p.NewPeerStream(listenPort)
	if err != nil {
		panic(err)
	}

	ps.Run(ctx, hostGroupName)
	select {}
}

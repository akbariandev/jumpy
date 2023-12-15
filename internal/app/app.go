package app

import (
	"context"
	"github.com/akbariandev/jumpy/internal/chain"
	"github.com/akbariandev/jumpy/internal/p2p"
	_ "github.com/libp2p/go-libp2p/p2p/host/peerstore"
)

const defaultHostGroupName = "jumpy"

func Start(listenPort int, hostGroupName string) {
	ctx := context.Background()

	//init genesis block
	genesisBlock := chain.CreateGenesisBlock("genesis_block_data")
	chain.Blockchain = append(chain.Blockchain, genesisBlock)

	if len(hostGroupName) == 0 {
		hostGroupName = defaultHostGroupName
	}

	ps := p2p.NewPeerStream(listenPort)
	ps.Run(ctx, hostGroupName)
	select {}
}

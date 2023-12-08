package app

import (
	"context"
	"github.com/akbariandev/jumpy/internal/chain"
	"github.com/akbariandev/jumpy/internal/p2p"
	_ "github.com/libp2p/go-libp2p/p2p/host/peerstore"
)

const hostGroupName = "jumpy"

func Start(listenPort int) {
	ctx := context.Background()

	//init genesis block
	genesisBlock := chain.CreateGenesisBlock("genesis_block_data")
	chain.Blockchain = append(chain.Blockchain, genesisBlock)

	p2p.Run(ctx, listenPort, hostGroupName)
	select {}
}

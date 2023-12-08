package app

import (
	"context"
	_ "github.com/libp2p/go-libp2p/p2p/host/peerstore"
	"myBlockchain/internal/chain"
	"myBlockchain/internal/p2p"
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

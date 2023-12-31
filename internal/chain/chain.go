package chain

import (
	"fmt"
	mrand "math/rand"
)

const (
	SuccessColor = "\033[1;32m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	InfoColor    = "\033[1;34m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

type Chain []Block

func (c Chain) GetLastBlock() *Block {
	i := len(c) - 1
	return &c[i]
}

func (c Chain) GetBlockAtIndex(index int) *Block {
	return &c[index]
}

func (c Chain) GetBlockByHash(hash string) *Block {
	for _, b := range c {
		if b.Hash == hash {
			return &b
		}
	}

	return nil
}

func (c Chain) GetRandomBlock() *Block {
	return &c[mrand.Intn(len(c))]
}

func (c Chain) PrintBlockChain() {
	for _, b := range c {
		fmt.Println("-----------------------------------------------------------------------")
		fmt.Println(fmt.Sprintf("Index: %d\nHash:%s", b.Index, b.Hash))
		fmt.Printf(SuccessColor, "Transactions:\n")
		for i, t := range b.Transaction {
			fmt.Println(fmt.Sprintf("%d: = %s", i, t.Data))
		}
		fmt.Printf(ErrorColor, "Conenctions:\n")
		for _, c := range b.Connections {
			fmt.Println(fmt.Sprintf("Node = %s", c.PeerID))
			fmt.Println(fmt.Sprintf("Block = %s", c.BlockHash))
		}
	}
}

func (c Chain) GetHeight() int {
	return len(c)
}

package chain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Transaction struct {
	Data any `json:"data"`
}

type BlockConnection struct {
	PeerID    string
	BlockHash string
}

type Block struct {
	Index       int
	Timestamp   string
	Transaction []Transaction
	Hash        string
	Connections []BlockConnection
}

func CreateGenesisBlock(data string) Block {
	t := time.Now()
	genesisBlock := Block{}
	connections := make([]BlockConnection, 2)
	return Block{0, t.String(), []Transaction{{Data: data}}, genesisBlock.calculateHash(), connections}
}

func (b Block) calculateHash() string {
	record := b.toString()
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func (b Block) toString() string {
	transactionStr := ""
	for _, t := range b.Transaction {
		transactionStr = fmt.Sprintf("%s%s", transactionStr, t.Data)
	}

	connectionStr := ""
	for _, c := range b.Connections {
		connectionStr = fmt.Sprintf("%s%s%s", connectionStr, c.PeerID, c.BlockHash)
	}

	return fmt.Sprintf("%d%s%s%s", b.Index, b.Timestamp, transactionStr, connectionStr)
}

// GenerateBlock will create a new block using previous block's hash
func GenerateBlock(localPeerID string, lastBlock *Block, targetBlockPeerID, targetBlockHash string, transaction []Transaction) Block {

	var newBlock Block

	newBlock.Index = lastBlock.Index + 1
	newBlock.Transaction = transaction
	connections := make([]BlockConnection, 0)
	connections = append(connections, BlockConnection{PeerID: localPeerID, BlockHash: lastBlock.Hash})
	connections = append(connections, BlockConnection{PeerID: targetBlockPeerID, BlockHash: targetBlockHash})
	newBlock.Connections = connections
	t := time.Now()
	newBlock.Timestamp = t.String()

	newBlock.Hash = newBlock.calculateHash()
	return newBlock
}

// GenerateMemoBlock will create a new block to store in memo before commitment
func GenerateMemoBlock(localPeerID string, lastBlock *Block, transaction []Transaction) *Block {

	newBlock := new(Block)

	newBlock.Index = lastBlock.Index + 1
	newBlock.Transaction = transaction
	connections := make([]BlockConnection, 0)
	connections = append(connections, BlockConnection{PeerID: localPeerID, BlockHash: lastBlock.Hash})
	newBlock.Connections = connections
	t := time.Now()
	newBlock.Timestamp = t.String()

	newBlock.Hash = newBlock.calculateHash()
	return newBlock
}

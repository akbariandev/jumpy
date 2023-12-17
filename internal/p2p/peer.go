package p2p

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/akbariandev/jumpy/internal/chain"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	net "github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"io"
	"log"
	mrand "math/rand"
	"strings"
)

const defaultBufSize = 4096

type PeerStream struct {
	Host            host.Host
	connections     map[string]*bufio.ReadWriter
	memTransactions []chain.Transaction
	memBlock        []*chain.Block
	chain           chain.Chain
	approveFunction func(block *chain.Block) bool
}

func NewPeerStream(listenPort int, approveFunc func(block *chain.Block) bool) (*PeerStream, error) {
	h, err := createHost(listenPort)
	if err != nil {
		return nil, err
	}

	ps := &PeerStream{
		Host:            h,
		memTransactions: make([]chain.Transaction, 0),
		connections:     make(map[string]*bufio.ReadWriter),
		memBlock:        make([]*chain.Block, 0),
		chain:           chain.Chain{},
		approveFunction: approveFunc,
	}

	//initialize genesis block
	genesisBlock := chain.CreateGenesisBlock("genesis_block_data")
	ps.chain = append(ps.chain, genesisBlock)

	return ps, nil
}

func (ps *PeerStream) Run(ctx context.Context, streamGroup string, runCli bool) {
	peerAddr := ps.getPeerFullAddr()
	log.Printf("my address: %s\n", peerAddr)

	if runCli {
		go ps.readCli()
	}

	// connect to other peers
	ps.Host.SetStreamHandler("/p2p/1.0.0", ps.handleStream)
	log.Println("listening for connections")
	peerChan := initMDNS(ps.Host, streamGroup)
	go func(ctx context.Context) {
		for {
			peer := <-peerChan
			if err := ps.Host.Connect(ctx, peer); err != nil {
				fmt.Println("connection failed:", err)
				continue
			}

			fmt.Println("connected to: ", peer)
			s, err := ps.Host.NewStream(ctx, peer.ID, "/p2p/1.0.0")
			if err != nil {
				fmt.Println("stream open failed", err)
			} else {
				rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
				ps.connections[peer.ID.String()] = rw
				go ps.readStream(s)
			}
		}
	}(ctx)
}

func (ps *PeerStream) ConnectionsIDs() []string {
	connections := []string{}
	for c, _ := range ps.connections {
		connections = append(connections, c)
	}

	return connections
}

func (ps *PeerStream) CommitTransaction() error {
	//generating memo block
	lastBlock := ps.chain.GetLastBlock()
	memBlock := chain.GenerateMemoBlock(ps.Host.ID().String(), lastBlock, ps.memTransactions)
	ps.memBlock = append(ps.memBlock, memBlock)
	ps.memTransactions = make([]chain.Transaction, 0)

	//send request approval message
	committee := ps.ConnectionsIDs()
	message := NewMessage(RequestApprovalTopic, RequestApprovalMessage{Block: memBlock, Committee: committee})
	if err := message.write(ps.connections[committee[0]]); err != nil {
		return err
	}
	return nil
}

func (ps *PeerStream) AddTransaction(data any) {
	ps.memTransactions = append(ps.memTransactions, chain.Transaction{
		Data: data,
	})
}

func (ps *PeerStream) GetChain() chain.Chain {
	return ps.chain
}

func (ps *PeerStream) handleStream(s net.Stream) {
	go ps.readStream(s)
}

func (ps *PeerStream) readStream(s net.Stream) {

	for {
		buffer := make([]byte, defaultBufSize)
		n, err := s.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from stream:", err)
			}
			break
		}
		b := buffer[:n]
		b = bytes.Trim(b, "\x00")
		m := &Message{}
		err = json.Unmarshal(b, m)
		if err != nil {
			fmt.Println(err)
			continue
		}
		switch m.Topic {
		case RequestTipBlockTopic:
			msg := &RequestTipBlockMessage{}
			if err = m.Payload.parse(msg); err != nil {
				fmt.Println(err)
				continue
			}
			if len(s.Conn().RemotePeer()) == 0 || ps.connections[s.Conn().RemotePeer().String()] == nil {
				fmt.Println(errors.New("sender ID is invalid"))
				continue
			}

			lastBlock := ps.chain.GetLastBlock()
			if lastBlock == nil {
				fmt.Println(errors.New("no block founded in chain"))
				continue
			}
			message := NewMessage(ResponseTipBlockTopic, ResponseTipBlockMessage{
				BlockHash: lastBlock.Hash,
			})

			if err = message.write(ps.connections[s.Conn().RemotePeer().String()]); err != nil {
				log.Println(err)
				continue
			}
		case ResponseTipBlockTopic:
			msg := &ResponseTipBlockMessage{}
			if err = m.Payload.parse(msg); err != nil {
				fmt.Println(err)
				continue
			}

			remotePeer := s.Conn().RemotePeer().String()
			if len(remotePeer) == 0 {
				fmt.Println(errors.New("sender ID is empty"))
				continue
			}
			if len(ps.memBlock) == 0 {
				fmt.Println(errors.New("no memo block to commit"))
				continue
			}

			block := ps.memBlock[0]
			block.Connections = append(block.Connections, chain.BlockConnection{PeerID: remotePeer, BlockHash: msg.BlockHash})
			ps.chain = append(ps.chain, *block)
			ps.memBlock = ps.memBlock[1:]
		case RequestApprovalTopic:
			msg := &RequestApprovalMessage{}
			if err = m.Payload.parse(msg); err != nil {
				fmt.Println(err)
				continue
			}

			resp := ResponseApprovalMessage{BlockHash: msg.Block.Hash, Committee: msg.Committee}
			resp.IsApproved = ps.approveFunction(msg.Block)
			message := NewMessage(ResponseApprovalTopic, resp)
			if err = message.write(ps.connections[s.Conn().RemotePeer().String()]); err != nil {
				log.Println(err)
				continue
			}
		case ResponseApprovalTopic:
			msg := &ResponseApprovalMessage{}
			if err = m.Payload.parse(msg); err != nil {
				fmt.Println(err)
				continue
			}
			if !msg.IsApproved {
				//failed approval
				//TODO: remove memBlock
				continue
			}

			msg.Committee = removePeerInCommittee(msg.Committee, s.Conn().RemotePeer().String())
			if len(msg.Committee) == 0 {
				//go for getting tip block
				randomPeerID := ps.getRandomPeer()
				if len(ps.connections) == 0 {
					//TODO: handle err
					continue
				}

				// create message and send to target peer
				message := NewMessage(RequestTipBlockTopic, RequestTipBlockMessage{})
				if err := message.write(ps.connections[randomPeerID.String()]); err != nil {
					//TODO: handle err
					continue
				}
				continue
			}

			memBlock := ps.getMemBlockByHash(msg.BlockHash)
			if memBlock == nil {
				//TODO: handle err
				continue
			}
			message := NewMessage(RequestApprovalTopic, RequestApprovalMessage{Block: memBlock, Committee: msg.Committee})
			if err := message.write(ps.connections[msg.Committee[0]]); err != nil {
				//TODO: handle err
				continue
			}
		default:
			fmt.Println(errors.New("undefined message"))
			continue
		}
	}
}

func (ps *PeerStream) getRandomPeer() peer.ID {
	var randomPeer peer.ID
	for {
		peersLen := ps.Host.Peerstore().Peers().Len()
		randomIndex := mrand.Intn(peersLen)
		randomPeer = ps.Host.Peerstore().Peers()[randomIndex]
		if randomPeer.String() == ps.Host.ID().String() {
			continue
		}
		break
	}

	return randomPeer
}

func (ps *PeerStream) getPeerFullAddr() ma.Multiaddr {
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", ps.Host.ID()))

	addrs := ps.Host.Addrs()
	var addr ma.Multiaddr
	for _, i := range addrs {
		if strings.HasPrefix(i.String(), "/ip4") {
			addr = i
			break
		}
	}
	return addr.Encapsulate(hostAddr)
}

func createHost(listenPort int) (host.Host, error) {

	r := rand.Reader
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	sourceMultiAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort))
	opts := []libp2p.Option{
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(priv),
	}

	host, err := libp2p.New(opts...)
	if err != nil {
		return nil, err
	}
	return host, nil
}

func removePeerInCommittee(in []string, peer string) (out []string) {
	for _, c := range in {
		if peer != c {
			out = append(out, c)
		}
	}

	return out
}

func (ps *PeerStream) getMemBlockByHash(hash string) *chain.Block {
	for _, b := range ps.memBlock {
		if b.Hash == hash {
			return b
		}
	}

	return nil
}

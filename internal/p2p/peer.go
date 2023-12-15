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
	MemTransactions []chain.Transaction
}

var connections map[string]*bufio.ReadWriter

func NewPeerStream(listenPort int) *PeerStream {
	h, err := createHost(listenPort)
	if err != nil {
		log.Fatal(err)
	}
	memTransactions := make([]chain.Transaction, 0)
	connections = make(map[string]*bufio.ReadWriter)

	return &PeerStream{Host: h, MemTransactions: memTransactions}
}

func (ps *PeerStream) handleStream(s net.Stream) {
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	go ps.readStream(s, rw)
}

func (ps *PeerStream) readStream(s net.Stream, rw *bufio.ReadWriter) {

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
		msg := &Message{}
		err = json.Unmarshal(b, msg)
		if err != nil {
			fmt.Println(err)
			continue
		}
		switch msg.Topic {
		case PullBlockTopic:
			pullMsg := &PullBlockMessage{}
			if err = msg.Payload.parse(pullMsg); err != nil {
				fmt.Println(err)
				continue
			}
			if len(pullMsg.SelfID) == 0 {
				fmt.Println(errors.New("sender ID is empty"))
				continue
			}

			lastBlock := chain.GetLastBlock()
			if lastBlock == nil {
				fmt.Println(errors.New("no block founded in chain"))
				continue
			}

			message := NewMessage(PushBlockTopic, PushBlockMessage{
				SelfID:    ps.Host.ID(),
				TargetID:  pullMsg.SelfID,
				BlockHash: lastBlock.Hash,
			})
			if err = message.write(rw); err != nil {
				log.Println(err)
				continue
			}
		case PushBlockTopic:
			pushMsg := &PushBlockMessage{}
			if err = msg.Payload.parse(pushMsg); err != nil {
				fmt.Println(err)
				continue
			}

			lastBlock := chain.GetLastBlock()
			if len(pushMsg.SelfID.String()) == 0 {
				fmt.Println(errors.New("sender ID is empty"))
				continue
			}

			if pushMsg.TargetID.String() == ps.Host.ID().String() {
				block := chain.GenerateBlock(ps.Host.ID().String(), lastBlock, pushMsg.TargetID.String(), pushMsg.BlockHash, ps.MemTransactions)
				chain.Blockchain = append(chain.Blockchain, block)
				ps.MemTransactions = make([]chain.Transaction, 0)
			} else {
				fmt.Println(errors.New("sender ID is not equal to my ID"))
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

func (ps *PeerStream) Run(ctx context.Context, streamGroup string) {
	peerAddr := ps.getPeerFullAddr()
	log.Printf("my address: %s\n", peerAddr)
	go ps.readCli()

	// connect to other peers
	ps.Host.SetStreamHandler("/p2p/1.0.0", ps.handleStream)
	log.Println("listening for connections")
	peerChan := initMDNS(ps.Host, streamGroup)
	go func(ctx context.Context, ps *PeerStream) {
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
				go ps.readStream(s, rw)
				connections[peer.ID.String()] = rw
			}
		}
	}(ctx, ps)
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

func addTransaction(ps *PeerStream, data any) {
	ps.MemTransactions = append(ps.MemTransactions, chain.Transaction{
		Data: data,
	})
}

func commitTransaction(ps *PeerStream) {
	randomPeerID := ps.getRandomPeer()
	message := NewMessage(PullBlockTopic, PullBlockMessage{
		SelfID: ps.Host.ID(),
	})
	if err := message.write(connections[randomPeerID.String()]); err != nil {
		fmt.Println(err)
	}
}

package p2p

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/akbariandev/jumpy/internal/chain"
)

type MessageTopic string
type MessagePayload []byte

const (
	RequestTipBlockTopic  MessageTopic = "request_tip_block"
	ResponseTipBlockTopic MessageTopic = "response_tip_block"
	RequestApprovalTopic  MessageTopic = "request_approve_block"
	ResponseApprovalTopic MessageTopic = "response_approve_block"
)

type Message struct {
	Topic   MessageTopic   `json:"t"`
	Payload MessagePayload `json:"p"`
}

type RequestTipBlockMessage struct{}

type ResponseTipBlockMessage struct {
	BlockHash string `json:"b"`
}

type RequestApprovalMessage struct {
	Block     *chain.Block `json:"b"`
	Committee []string     `json:"c"`
}

type ResponseApprovalMessage struct {
	BlockHash  string   `json:"h"`
	IsApproved bool     `json:"a"`
	Committee  []string `json:"c"`
}

func NewMessage(topic MessageTopic, payload any) *Message {
	pByte, err := json.Marshal(payload)
	if err != nil {
		return nil
	}
	return &Message{
		Topic:   topic,
		Payload: pByte,
	}
}

func (m *Message) write(rw *bufio.ReadWriter) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if len(b) > defaultBufSize {
		return errors.New("message size exceeded")
	}
	bb := make([]byte, 0)
	padding := make([]byte, defaultBufSize-len(b))
	bb = append(b, padding...)
	if _, err := rw.Write(bb); err != nil {
		return err
	}
	return rw.Flush()
}

func (p MessagePayload) parse(msg any) error {
	if err := json.Unmarshal(p, msg); err != nil {
		return err
	}
	return nil
}

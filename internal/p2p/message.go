package p2p

import (
	"bufio"
	"encoding/json"
	"errors"
)

type MessageTopic string
type MessagePayload []byte

type Message struct {
	Topic   MessageTopic   `json:"topic"`
	Payload MessagePayload `json:"payload"`
}

type PullBlockMessage struct{}

type PushBlockMessage struct {
	BlockHash string `json:"b"`
}

const (
	PullBlockTopic MessageTopic = "pull_block"
	PushBlockTopic MessageTopic = "push_block"
)

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

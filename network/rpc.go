/***************************************************************
 * Arquivo: rpc.go
 * Descrição: Implementação de RPC para comunicação de rede.
 * Autor: JoaoRafa19
 * Data de criação: 2024-2025
 * Versão:
 * Licença:
 * Observações:
 ***************************************************************/

package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"

	"github.com/JoaoRafa19/crypto-go/core"
	"github.com/sirupsen/logrus"
)

type MessageType byte

const (
	MessageTypeTx MessageType = 0x0
	MessageTypeBlock
)

type Message struct {
	Header MessageType
	Data   []byte
}

func NewMessage(t MessageType, data []byte) *Message {
	return &Message{
		Header: t,
		Data:   data,
	}
}

func (msg *Message) Bytes() []byte {
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(msg)
	return buf.Bytes()
}

type RPC struct {
	From    NetAddr
	Payload io.Reader
}

type DecodedMessage struct {
	From NetAddr
	Data any
}
type RPCDecodeFunc func(RPC) (*DecodedMessage, error)

func DefaultRPCDecodeFunc(rpc RPC) (*DecodedMessage, error) {

	msg := Message{}
	if err := gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil {
		return nil, fmt.Errorf("failed to decode message from %s:%s", rpc.From, err)
	}
	logrus.WithFields(logrus.Fields{
		"from": rpc.From,
		"type": msg.Header,
	}).Debug("new incomming message")
	
	switch msg.Header {
	case MessageTypeTx:
		tx := new(core.Transaction)
		if err := tx.Decode(core.NewGobDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, err
		}
		return &DecodedMessage{From: rpc.From, Data: tx}, nil
	default:
		return nil, fmt.Errorf("invalid message header %d", msg.Header)
	}
}

type RPCProcessor interface {
	ProcessMessage(*DecodedMessage) error
}

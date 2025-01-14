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
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io"

	"github.com/JoaoRafa19/crypto-go/core"
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

type RPCHandler interface {
	HandleRPC(rpc RPC) error
}

type DefaultRPCHandler struct {
	P RPCProcessor
}

func NewDefaultRPCHandler(p RPCProcessor) RPCHandler {
	gob.Register(elliptic.P256())

	return &DefaultRPCHandler{
		P: p,
	}
}

func (h *DefaultRPCHandler) HandleRPC(rpc RPC) error {
	msg := Message{}
	if err := gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil {
		return fmt.Errorf("failed to decode message from %s:%s", rpc.From, err)
	}
	switch msg.Header {
	case MessageTypeTx:
		tx := new(core.Transaction)
		if err := tx.Decode(core.NewGobDecoder(bytes.NewReader(msg.Data))); err != nil {
			return err
		}
		return h.P.ProcessTransaction(rpc.From, tx)
	default:
		return fmt.Errorf("invalid message header %d", msg.Header)
	}
}

type RPCProcessor interface {
	ProcessTransaction(NetAddr, *core.Transaction) error
}

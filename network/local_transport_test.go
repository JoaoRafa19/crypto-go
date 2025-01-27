/***************************************************************
 * Arquivo: local_transport_test.go
 * Descrição: Teste do local transport
 * Autor: JoaoRafa19
 * Data de criação: 2024-2025
 * Versão: 0.0.1
 * Licença: MIT License
 * Observações:
 ***************************************************************/
package network

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestConnect(t *testing.T) {
	tra := NewLocalTransport("A").(*LocalTransport)
	trb := NewLocalTransport("B").(*LocalTransport)

	tra.Connect(trb)
	trb.Connect(tra)

	assert.Equal(t, tra.Peers[trb.addr], trb)
	assert.Equal(t, trb.Peers[tra.addr], tra)
}

func TestMessage(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")

	tra.Connect(trb)
	trb.Connect(tra)
	msg := []byte("hello world")
	assert.Nil(t, tra.SendMessage(trb.Addr(), msg))

	rpc := <-trb.Consume()
	b, err := io.ReadAll(rpc.Payload)

	assert.Nil(t, err)
	assert.Equal(t, b, msg)
	assert.Equal(t, rpc.From, tra.Addr())
}

func TestBroadcast(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")
	trc := NewLocalTransport("C")

	tra.Connect(trb)
	tra.Connect(trc)

	msg := []byte("hello world")
	assert.Nil(t, tra.Broadcast(msg))

	rpcB := <-trb.Consume()
	b, err := io.ReadAll(rpcB.Payload)
	assert.Nil(t, err)
	assert.Equal(t, b, msg)

	rpcC := <-trc.Consume()
	c, err := io.ReadAll(rpcC.Payload)
	assert.Nil(t, err)
	assert.Equal(t, c, msg)
}

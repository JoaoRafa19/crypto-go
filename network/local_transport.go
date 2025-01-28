/***************************************************************
 * Arquivo: local_transport.go
 * Descrição: Implementação do transporte local.
 * Autor: JoaoRafa19
 * Data de criação: 2024-2025
 * Versão: 0.0.1
 * Licença: MIT License
 * Observações:
 ***************************************************************/

package network

import (
	"bytes"
	"fmt"
	"sync"
)

type LocalTransport struct {
	addr      NetAddr
	Peers     map[NetAddr]*LocalTransport
	Lock      sync.RWMutex
	ConsumeCh chan RPC
}

func NewLocalTransport(addr NetAddr) Transport {
	return &LocalTransport{
		addr:      addr,
		ConsumeCh: make(chan RPC, 1024),
		Peers:     make(map[NetAddr]*LocalTransport),
	}
}

func (t *LocalTransport) Consume() <-chan RPC {
	return t.ConsumeCh

}
func (t *LocalTransport) Connect(tr Transport) error {
	t.Lock.Lock()
	defer t.Lock.Unlock()

	t.Peers[tr.Addr()] = tr.(*LocalTransport)

	return nil
}

func (t *LocalTransport) SendMessage(to NetAddr, payload []byte) error {
	t.Lock.RLock()
	defer t.Lock.RUnlock()

	perr, ok := t.Peers[to]
	if !ok {
		return fmt.Errorf("%s could not send message to unkown peer %s", t.addr, to)

	}

	perr.ConsumeCh <- RPC{
		From:    t.addr,
		Payload: bytes.NewReader(payload),
	}
	return nil
}

func (t *LocalTransport) Addr() NetAddr {
	return t.addr
}

func (t *LocalTransport) Broadcast(payload []byte) error {
	for _, peer := range t.Peers {
		if err := t.SendMessage(peer.Addr(), payload); err != nil {
			return err
		}
	}
	return nil
}

/***************************************************************
 * Arquivo: main.go
 * Descrição: Ponto de entrada da aplicação.
 * Autor: JoaoRafa19
 * Data de criação: 2024-2025
 * Versão: 0.0.1
 * Licença: MIT License
 * Observações:
 ***************************************************************/

package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"strconv"
	"time"

	"github.com/JoaoRafa19/crypto-go/core"
	"github.com/JoaoRafa19/crypto-go/crypto"
	"github.com/JoaoRafa19/crypto-go/network"
	"github.com/JoaoRafa19/crypto-go/types"
	"github.com/sirupsen/logrus"
)

func init() {
	gob.Register(&big.Int{})
	gob.Register(&ecdsa.PublicKey{})
	gob.Register(&crypto.PublicKey{})
	gob.Register(&crypto.Signature{})
	gob.Register(&types.Hash{})
}

func main() {
	trLocal := network.NewLocalTransport("LOCAL")
	trRemoteA := network.NewLocalTransport("REMOTE_A") // 24.123.123
	trRemoteB := network.NewLocalTransport("REMOTE_B")
	trRemoteC := network.NewLocalTransport("REMOTE_C")

	trLocal.Connect(trRemoteA)
	trRemoteA.Connect(trRemoteB)
	trRemoteB.Connect(trRemoteC)

	trRemoteA.Connect(trLocal)

	initRemoteServers([]network.Transport{trRemoteA, trRemoteB, trRemoteC})

	go func() {
		for {
			// trRemote.SendMessage(trLocal.Addr(), []byte("hello world"))
			if err := sendTransaction(trRemoteA, trLocal.Addr()); err != nil {
				logrus.Error(err)
			}
			time.Sleep(time.Second * 2)
		}
	}()

	privKey := crypto.GeneratePrivateKey()

	localServer := makeServer("LOCAL", trLocal, &privKey)
	localServer.Start()

}

func initRemoteServers(trs []network.Transport) {
	for index, transport := range trs {
		id := fmt.Sprintf("REMOTE_%d", index)
		s := makeServer(id, transport, nil)
		go s.Start()
	}
}

func makeServer(id string, tr network.Transport, privKey *crypto.PrivateKey) *network.Server {
	opts := network.ServerOpts{
		Transports: []network.Transport{tr},
		PrivateKey: privKey,
		ID:         id,
	}
	s, err := network.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func sendTransaction(tr network.Transport, to network.NetAddr) error {
	privKey := crypto.GeneratePrivateKey()

	data := []byte(strconv.FormatInt(int64(rand.Intn(10000000000000)), 10))
	tx := core.NewTransaction(data)
	tx.Sign(privKey)
	buf := &bytes.Buffer{}

	if err := tx.Encode(core.NewGobEncoder(buf)); err != nil {
		return err
	}

	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())
	return tr.SendMessage(to, msg.Bytes())

}

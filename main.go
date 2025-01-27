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
	trRemote := network.NewLocalTransport("REMOTE") // 24.123.123

	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	go func() {
		for {
			// trRemote.SendMessage(trLocal.Addr(), []byte("hello world"))
			if err := sendTransaction(trRemote, trLocal.Addr()); err != nil {
				logrus.Error(err)
			}
			time.Sleep(time.Second * 2)
		}
	}()
	privKey := crypto.GeneratePrivateKey()

	opts := network.ServerOpts{
		Transports: []network.Transport{
			trLocal,
		},
		PrivateKey: &privKey,
		ID:         "LOCAL",
	}

	s, err := network.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}
	s.Start()
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

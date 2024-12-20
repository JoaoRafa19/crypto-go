package main

import (
	"time"

	"github.com/JoaoRafa19/crypto-go/network"
)

func main() {
	trLocal := network.NewLocalTransport("LOCAL")
	trRemote := network.NewLocalTransport("REMOTE")

	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	go func() {
		for {
			trRemote.SendMessage(trLocal.Addr(), []byte("hello world"))
			time.Sleep(time.Second * 1)
		}
	}()

	opts := network.ServerOpts{
		Transports: []network.Transport{
			trLocal,
		},
		
	}

	s := network.NewServer(opts)

	s.Start()
}

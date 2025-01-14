/***************************************************************
 * Arquivo: server.go
 * Descrição: Implementação do servidor de rede.
 * Autor: JoaoRafa19
 * Data de criação: 2024-2025
 * Versão: 0.0.1
 * Licença: MIT License
 * Observações: 
 ***************************************************************/

package network

import (
	"fmt"
	"time"

	"github.com/JoaoRafa19/crypto-go/core"
	"github.com/JoaoRafa19/crypto-go/crypto"
	"github.com/sirupsen/logrus"
)

var defaultBlockTime = time.Duration(time.Second * 5)

type ServerOpts struct {
	// RPCHandler is responsible for handling remote procedure calls (RPCs)
	// within the network server. It defines the methods and logic required
	// to process incoming RPC requests and send appropriate responses.
	RPCHandler RPCHandler
	Transports []Transport
	PrivateKey *crypto.PrivateKey
	BlockTime  time.Duration
}

type Server struct {
	ServerOpts
	MemPool     *TxPool
	IsValidator bool
	RpcCh       chan RPC
	QuitChan    chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}
	
	s := &Server{
		opts,
		NewTxPool(),
		opts.PrivateKey != nil,
		make(chan RPC),
		make(chan struct{}),
	}
	if opts.RPCHandler == nil {
		opts.RPCHandler = NewDefaultRPCHandler(s)
	}
	s.ServerOpts = opts
	return s
}

func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(s.BlockTime)
free:
	for {
		select {
		case rpc := <-s.RpcCh:
			fmt.Printf("%+v\n", rpc)
			if err := s.RPCHandler.HandleRPC(rpc); err != nil {
				logrus.Error(err)
			}
		case <-s.QuitChan:
			break free
		case <-ticker.C:
			if s.IsValidator {
				s.CreateNewBlock()
				fmt.Println("creating a new block")
			}
		}
	}
	fmt.Println("Server shutdown")
}

func (s *Server) ProcessTransaction(from NetAddr, tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})

	if s.MemPool.Contains(hash) {
		logrus.WithField(
			"Adding New tx to mempool",
			logrus.Fields{
				"hash": tx.Hash(core.TxHasher{}),
			},
		).Info("Transaction already in mempool")
		return nil
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	tx.SetFirstSeen(time.Now().UnixNano())

	logrus.WithField(
		"Adding New tx to mempool",
		logrus.Fields{
			"hash": hash,
		},
	).Info("Add to mempool")

	return s.MemPool.Add(tx)
}
func (s *Server) CreateNewBlock() error {
	fmt.Println("create a new block")
	return nil
}

func (s *Server) initTransports() {
	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				// handle
				s.RpcCh <- rpc
			}
		}(tr)
	}
}

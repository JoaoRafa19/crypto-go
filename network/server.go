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
	"bytes"
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
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	Transports    []Transport
	PrivateKey    *crypto.PrivateKey
	BlockTime     time.Duration
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
	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}

	s := &Server{
		opts,
		NewTxPool(),
		opts.PrivateKey != nil,
		make(chan RPC),
		make(chan struct{}),
	}

	s.ServerOpts = opts

	// if RPCProcessor is not provided, use the server as the
	// default RPC processor
	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}
	return s
}

func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(s.BlockTime)
free:
	for {
		select {
		case rpc := <-s.RpcCh:
			message, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				logrus.Warn(err)
				continue
			}
			if err := s.RPCProcessor.ProcessMessage(message); err != nil {
				logrus.Error(err)
				continue
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

func (s *Server) ProcessMessage(message *DecodedMessage) error {

	switch msg := message.Data.(type) {
	case *core.Transaction:
		return s.processTransaction(msg)
	default:
		return fmt.Errorf("unknown message type: %T", msg)
	}
}

func (s *Server) broadcast(payload []byte) error {
	for _, tr := range s.Transports {
		if err := tr.Broadcast(payload); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) processTransaction(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})

	if s.MemPool.Contains(hash) {
		logrus.WithFields(
			logrus.Fields{
				"hash":    fmt.Sprintf("%x", hash),
				"message": "Transaction already in mempool",
			},
		).Info("Transaction already in mempool")
		return nil
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	tx.SetFirstSeen(time.Now().UnixNano())

	logrus.WithFields(
		logrus.Fields{
			"hash":       fmt.Sprintf("%x", hash),
			"mempoollen": s.MemPool.Len(),
		},
	).Info("adding a new transaction to mempool\n")

	go s.broadcastTx(tx)

	return s.MemPool.Add(tx)
}
func (s *Server) CreateNewBlock() error {
	fmt.Println("create a new block")
	return nil
}

func (s *Server) broadcastTx(tx *core.Transaction) error {
	buf := &bytes.Buffer{}

	if err := tx.Encode(core.NewGobEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeTx, buf.Bytes())
	return s.broadcast(msg.Bytes())
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

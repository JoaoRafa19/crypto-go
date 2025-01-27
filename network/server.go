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
	"os"
	"time"

	"github.com/JoaoRafa19/crypto-go/core"
	"github.com/JoaoRafa19/crypto-go/crypto"
	"github.com/JoaoRafa19/crypto-go/types"
	"github.com/go-kit/log"
)

var defaultBlockTime = time.Duration(time.Second * 5)

type ServerOpts struct {
	ID     string
	Logger log.Logger
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
	chain       *core.BlockChain
	QuitChan    chan struct{}
}

func NewServer(opts ServerOpts) (*Server, error) {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}
	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}

	if opts.Logger == nil {
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		opts.Logger = log.With(opts.Logger, "ID", opts.ID)
	}
	chain, err := core.NewBlockChain(genesisBlock())
	if err != nil {
		return nil, err
	}

	s := &Server{
		ServerOpts:  opts,
		MemPool:     NewTxPool(),
		IsValidator: opts.PrivateKey != nil,
		RpcCh:       make(chan RPC),
		QuitChan:    make(chan struct{}),
		chain:       chain,
	}

	s.ServerOpts = opts

	// if RPCProcessor is not provided, use the server as the
	// default RPC processor
	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}

	if s.IsValidator {
		go s.ValidatorLoop()
	}
	return s, nil
}

func (s *Server) Start() {
	s.initTransports()

free:
	for {
		select {
		case rpc := <-s.RpcCh:
			message, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				s.Logger.Log("error", err)
				continue
			}
			if err := s.RPCProcessor.ProcessMessage(message); err != nil {
				s.Logger.Log("error", err)

				continue
			}

		case <-s.QuitChan:
			break free

		}
	}
	s.Logger.Log("msg", "Server shutdown")
}

func (s *Server) ValidatorLoop() {
	ticker := time.NewTicker(s.BlockTime)
	s.Logger.Log("msg", "Server starting validate", "block time: ", s.BlockTime)

	for {
		<-ticker.C
		s.CreateNewBlock()
	}
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
		return nil
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	tx.SetFirstSeen(time.Now().UnixNano())

	s.Logger.Log(
		"msg", "add transaction to mempool",
		"hash", hash,
		"mempoollen", s.MemPool.Len(),
	)

	go s.broadcastTx(tx)

	return s.MemPool.Add(tx)
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
func (s *Server) CreateNewBlock() error {
	currentHeader, err := s.chain.GetHeader(s.chain.Height())
	if err != nil {
		return err
	}

	block, err := core.NewBlockFromHeader(currentHeader, nil)
	if err != nil {
		return err
	}

	if err := block.Sign(*s.PrivateKey); err != nil {
		return err
	}

	if err := s.chain.AddBlock(block); err != nil {
		return err
	}

	return nil
}

func genesisBlock() *core.Block {
	header := &core.Header{
		Version:   1,
		DataHash:  types.Hash{},
		Timestamp: uint64(time.Now().UnixNano()),
		Height:    0,
	}

	b, _ := core.NewBlock(header, nil)
	return b
}

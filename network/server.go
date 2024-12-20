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
	Transports []Transport
	PrivateKey *crypto.PrivateKey
	BlockTime  time.Duration
}

type Server struct {
	ServerOpts
	memPool     *TxPool
	isValidator bool
	rpcCh       chan RPC
	blockTime   time.Duration
	quitChan    chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}
	return &Server{
		ServerOpts:  opts,
		memPool:     NewTxPool(),
		blockTime:   opts.BlockTime,
		isValidator: opts.PrivateKey != nil,
		rpcCh:       make(chan RPC),
		quitChan:    make(chan struct{}),
	}
}

func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(s.BlockTime)
free:
	for {
		select {
		case rpc := <-s.rpcCh:
			fmt.Printf("%+v\n", rpc)
		case <-s.quitChan:
			break free
		case <-ticker.C:
			if s.isValidator {
				s.CreateNewBlock()
				fmt.Println("creating a new block")
			}
		}
	}
	fmt.Println("Server shutdown")
}

func (s *Server) handleTransactions(tx *core.Transaction) error {
	if err := tx.Verify(); err != nil {
		return err
	}
	hash := tx.Hash(core.TxHasher{})

	if s.memPool.Has(hash) {
		logrus.WithField(
			"Adding New tx to mempool",
			logrus.Fields{
				"hash": tx.Hash(core.TxHasher{}),
			},
		).Info("Transaction already in mempool")
		return nil
	}

	logrus.WithField(
		"Adding New tx to mempool",
		logrus.Fields{
			"hash": hash,
		},
	).Info("Add to mempool")

	return s.memPool.Add(tx)
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
				s.rpcCh <- rpc
			}
		}(tr)
	}
}

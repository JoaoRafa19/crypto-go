package network

import (
	"sort"

	"github.com/JoaoRafa19/crypto-go/core"
	"github.com/JoaoRafa19/crypto-go/types"
)

type TxMapSorter struct {
	transations []*core.Transaction
}

func NewTxMapSorter(txMap map[types.Hash]*core.Transaction) *TxMapSorter {
	txx := make([]*core.Transaction, len(txMap))
	i := 0
	for _, tx := range txMap {
		txx[i] = tx
		i++
	}

	s := &TxMapSorter{transations: txx}

	sort.Sort(s)

	return s
}

func (s *TxMapSorter) Len() int {
	return len(s.transations)
}

// Swap
func (s *TxMapSorter) Swap(i, j int) {
	s.transations[i], s.transations[j] = s.transations[j], s.transations[i]
}

// Less
func (s *TxMapSorter) Less(i, j int) bool {
	return s.transations[i].FirstSeen() < s.transations[j].FirstSeen()
}

type TxPool struct {
	transactions map[types.Hash]*core.Transaction
}

func NewTxPool() *TxPool {
	return &TxPool{
		transactions: make(map[types.Hash]*core.Transaction),
	}
}

// Transactions returns a slice of all transactions in the pool.
func (p *TxPool) Transactions() []*core.Transaction {
	s := NewTxMapSorter(p.transactions)
	return s.transations
}

// Add adds a transaction to the pool. The caller is responsible for checking
// if the transaction already exists in the pool.
func (p *TxPool) Add(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})
	if p.Has(hash) {
		return nil
	}
	p.transactions[hash] = tx
	return nil
}

// Has checks if a transaction with the given hash exists in the pool.
func (p *TxPool) Has(hash types.Hash) bool {
	_, ok := p.transactions[hash]
	return ok
}

// Len returns the number of transactions currently in the pool.
func (p *TxPool) Len() int {
	return len(p.transactions)
}

// Flush removes all transactions from the pool.
func (p *TxPool) Flush() {
	p.transactions = make(map[types.Hash]*core.Transaction)
}

/***************************************************************
 * Arquivo: txpool.go
 * Descrição: Implementação do pool de transações.
 * Autor: JoaoRafa19
 * Data de criação: 2024-2025
 * Versão: 0.0.1
 * Licença: MIT License
 * Observações:
 ***************************************************************/

package network

import (
	"sort"

	"github.com/JoaoRafa19/crypto-go/core"
	"github.com/JoaoRafa19/crypto-go/types"
)

type TxMapSorter struct {
	Transations []*core.Transaction
}

func NewTxMapSorter(txMap map[types.Hash]*core.Transaction) *TxMapSorter {
	txx := make([]*core.Transaction, len(txMap))
	i := 0
	for _, tx := range txMap {
		txx[i] = tx
		i++
	}

	s := &TxMapSorter{Transations: txx}

	sort.Sort(s)

	return s
}

func (s *TxMapSorter) Len() int {
	return len(s.Transations)
}

// Swap
func (s *TxMapSorter) Swap(i, j int) {
	s.Transations[i], s.Transations[j] = s.Transations[j], s.Transations[i]
}

// Less
func (s *TxMapSorter) Less(i, j int) bool {
	return s.Transations[i].GetFirstSeen() < s.Transations[j].GetFirstSeen()
}

type TxPool struct {
	trxs map[types.Hash]*core.Transaction
}

func NewTxPool() *TxPool {
	return &TxPool{
		trxs: make(map[types.Hash]*core.Transaction),
	}
}

// Transactions returns a slice of all transactions in the pool.
func (p *TxPool) Transactions() []*core.Transaction {
	s := NewTxMapSorter(p.trxs)
	return s.Transations
}

// Add adds a transaction to the pool. The caller is responsible for checking
// if the transaction already exists in the pool.
func (p *TxPool) Add(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})
	if p.Contains(hash) {
		return nil
	}
	p.trxs[hash] = tx
	return nil
}

// Has checks if a transaction with the given hash exists in the pool.
func (p *TxPool) Contains(hash types.Hash) bool {
	_, ok := p.trxs[hash]
	return ok
}

// Len returns the number of transactions currently in the pool.
func (p *TxPool) Len() int {
	return len(p.trxs)
}

// Flush removes all transactions from the pool.
func (p *TxPool) Flush() {
	p.trxs = make(map[types.Hash]*core.Transaction)
}

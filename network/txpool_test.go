package network

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/JoaoRafa19/crypto-go/core"
	"github.com/stretchr/testify/assert"
)

func TestTxPool(t *testing.T) {
	p := NewTxPool()
	assert.Equal(t, p.Len(), 0)
}

func TestTxPoolAddTx(t *testing.T) {
	p := NewTxPool()
	tx := core.NewTransaction([]byte("foo bar baz"))
	assert.Nil(t, p.Add(tx))
	assert.Equal(t, p.Len(), 1)

	_ = core.NewTransaction([]byte("foo"))
	assert.Equal(t, p.Len(), 1)

	p.Flush()
	assert.Empty(t, p.trxs)
	assert.Zero(t, p.Len())

}

func TestSortTransactions(t *testing.T) {
	p := NewTxPool()

	txlen := 1000
	for i := 0; i < txlen; i++ {
		tx := core.NewTransaction([]byte(strconv.FormatInt(int64(i), 10)))
		// Gera um valor exclusivo e suficientemente aleatório para FirstSeen
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		randomDelay := int64(r.Intn(1000) + 1) // Adiciona uma variação aleatória e evita 0
		tx.SetFirstSeen(time.Now().UnixNano() * randomDelay * int64(i+1)) // evita 0
		
		assert.Nil(t, p.Add(tx))
	}

	assert.Equal(t, txlen, p.Len())

	txx := p.Transactions()

	for i := 0; i < len(txx)-1; i++ {
		
		x := txx[i].GetFirstSeen()
		x1 := txx[i+1].GetFirstSeen()
		assert.True(t, x < x1)
	}
}

/***************************************************************
 * Arquivo: blockchain.go
 * Descrição: Implementação da estrutura de blockchain.
 * Autor: JoaoRafa19
 * Data de criação: 2024-2025
 * Versão: 0.0.1
 * Licença: MIT License
 * Observações:
 ***************************************************************/

package core

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

type BlockChain struct {
	Store     Storage
	Lock      sync.RWMutex
	Headers   []*Header
	Validator Validator
}

func NewBlockChain(genesis *Block) (*BlockChain, error) {
	bc := &BlockChain{
		Headers: []*Header{},
		Store:   NewMemStore(),
	}
	bc.Validator = NewBlockValidator(bc)
	err := bc.addBlockWithoutValidation(genesis)

	return bc, err

}

func (bc *BlockChain) SetValidator(v Validator) {
	bc.Validator = v
}
func (bc *BlockChain) AddBlock(b *Block) error {
	//validate
	err := bc.Validator.ValidateBlock(b)
	if err != nil {
		return err
	}

	return bc.addBlockWithoutValidation(b)
}

func (bc *BlockChain) HasBlock(heigh uint32) bool {
	return heigh <= bc.Height()
}

// [g, 1, 2, 3] = len 4 ; heigh = 3
func (bc *BlockChain) Height() uint32 {
	bc.Lock.RLock()
	defer bc.Lock.RUnlock()
	return uint32(len(bc.Headers) - 1)
}

func (bc *BlockChain) addBlockWithoutValidation(b *Block) error {
	bc.Lock.Lock()
	bc.Headers = append(bc.Headers, b.Header)
	bc.Lock.Unlock()

	logrus.WithField("Adding New Block", logrus.Fields{
		"height": b.Height,
		"hash":   b.Hash(BlockHasher{}),
	}).Info("adding new block")

	return bc.Store.Put(b)
}

func (bc *BlockChain) GetHeader(height uint32) (*Header, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("height (%+v) is too high", height)
	}
	bc.Lock.Lock()
	defer bc.Lock.Unlock()

	return bc.Headers[height], nil
}

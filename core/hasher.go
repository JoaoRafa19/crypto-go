/***************************************************************
 * Arquivo: hasher.go
 * Descrição: Implementação da interface de hash.
 * Autor: JoaoRafa19
 * Data de criação: 2024-2025
 * Versão: 0.0.1
 * Licença: MIT License
 * Observações: 
 ***************************************************************/

package core

import (
	"crypto/sha256"

	"github.com/JoaoRafa19/crypto-go/types"
)

type Hasher[T any] interface {
	Hash(T) types.Hash
}

type BlockHasher struct{}

func (BlockHasher) Hash(b *Header) types.Hash {

	h := sha256.Sum256(b.Bytes())
	return types.Hash(h)
}

type TxHasher struct {

}

func (TxHasher) Hash(tx*Transaction) types.Hash {
	return types.Hash(sha256.Sum256(tx.Data)) 
}
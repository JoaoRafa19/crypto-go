/***************************************************************
 * Arquivo: storage.go
 * Descrição: Implementação da interface de armazenamento.
 * Autor: JoaoRafa19
 * Data de criação: 2024-2025
 * Versão: 0.0.1
 * Licença: MIT License
 * Observações: 
 ***************************************************************/

package core

type Storage interface {
	Put(*Block) error
}

type MemoryStore struct{}

func NewMemStore() *MemoryStore {
	return &MemoryStore{}
}

func (ms *MemoryStore) Put(b *Block) error {
	return nil
}

/***************************************************************
 * Arquivo: block.go
 * Descrição: Implementação da estrutura de bloco.
 * Autor: JoaoRafa19
 * Data de criação: 2024-2025
 * Versão: 0.0.1
 * Licença: MIT License
 * Observações:
 ***************************************************************/

package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/JoaoRafa19/crypto-go/crypto"
	"github.com/JoaoRafa19/crypto-go/types"
)

// To save space we just hash the header
type Header struct {
	Version       uint32
	PrevBlockHash types.Hash
	Timestamp     uint64
	Height        uint32
	DataHash      types.Hash
}

func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}

	enc := gob.NewEncoder(buf)
	enc.Encode(h)
	return buf.Bytes()
}

// Hold the transactions and the header information
type Block struct {
	*Header
	Transactions []*Transaction
	Validator    crypto.PublicKey
	Signature    *crypto.Signature

	//nonce uint32
	// cached version of the header hash
	hash types.Hash
}

func NewBlock(h *Header, tsc []*Transaction) (*Block, error ) {
	return &Block{
		Header:       h,
		Transactions: tsc,
	}, nil
}

func NewBlockFromHeader(prevHeader *Header, txx []*Transaction) (*Block, error) {
	dataHash, err := CalculateDataHash(txx)
	if err != nil {
		return nil, err
	}

	header := &Header{
		Version:       1,
		Height:        prevHeader.Height + 1,
		DataHash:      dataHash,
		PrevBlockHash: BlockHasher{}.Hash(prevHeader),
		Timestamp:     uint64(time.Now().UnixNano()),
	}

	return NewBlock(header, txx)
}

func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, tx)
}

func (b *Block) Verify() error {
	if b.Signature == nil {
		return fmt.Errorf("block has no signature")
	}

	if !b.Signature.Verify(b.Validator, b.Header.Bytes()) {
		return fmt.Errorf("block has invalid signature")
	}

	for _, trx := range b.Transactions {
		if err := trx.Verify(); err != nil {
			return err
		}
	}

	hash, err := CalculateDataHash(b.Transactions)
	if err != nil {
		return err
	}

	if hash != b.DataHash {
		return fmt.Errorf("block %s has an invalid data hash", b.Hash(BlockHasher{}))
	}

	return nil
}

func (b *Block) Sign(privKey crypto.PrivateKey) error {
	sig, err := privKey.Sign(b.Header.Bytes())
	if err != nil {
		return err
	}
	b.Validator = privKey.PublicKey()
	b.Signature = sig
	return nil
}

func (b *Block) Decode(dec Decoder[*Block]) error {
	return dec.Decode(b)
}

func (b *Block) Encode(enc Encoder[*Block]) error {
	return enc.Encode(b)
}

func (b *Block) Hash(hasher Hasher[*Header]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b.Header)
	}

	return b.hash
}

func CalculateDataHash(txx []*Transaction) (hash types.Hash, err error) {
	var (
		buf = &bytes.Buffer{}
	)

	for _, tx := range txx {
		if err = tx.Encode(NewGobEncoder(buf)); err != nil {
			return
		}
	}

	hash = sha256.Sum256(buf.Bytes())
	return
}

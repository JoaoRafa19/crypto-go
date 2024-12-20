package core

import (
	"fmt"

	"github.com/JoaoRafa19/crypto-go/crypto"
	"github.com/JoaoRafa19/crypto-go/types"
)

type Transaction struct {
	Data []byte

	From      crypto.PublicKey
	Signature *crypto.Signature

	//cached version of tx data hash
	CacheHash types.Hash
	// firstSeen is the tmiestamp of when this tx is first seen localy
	FirstSeen int64
}

func (tx *Transaction) Sign(privKey crypto.PrivateKey) error {
	sig, err := privKey.Sign(tx.Data)
	if err != nil {
		return err
	}

	tx.From = privKey.PublicKey()
	tx.Signature = sig
	return nil
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data: data,
	}
}

func (tx *Transaction) Hash(h Hasher[*Transaction]) types.Hash {
	if tx.CacheHash.IsZero() {
		tx.CacheHash = h.Hash(tx)
	}
	return h.Hash(tx)
}

func (tx *Transaction) Verify() error {
	if tx.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}

	if !tx.Signature.Verify(tx.From, tx.Data) {
		return fmt.Errorf("invalid transaction signature")
	}

	return nil
}

func (tx *Transaction) Decode(dec Decoder[*Transaction]) error {
	return dec.Decode(tx)
}
func (tx *Transaction) Encode(enc Encoder[*Transaction]) error {
	return enc.Encode(tx)
}

// set firstSeen
func (tx *Transaction) SetFirstSeen(t int64) {
	tx.FirstSeen = t
}

// get firstSeen
func (tx *Transaction) GetFirstSeen() int64 {
	return tx.FirstSeen
}

package core

import (
	"crypto/elliptic"
	"encoding/gob"
	"io"

	"github.com/JoaoRafa19/crypto-go/crypto"
)

type Encoder[T any] interface {
	Encode(T) error
}

type Decoder[T any] interface {
	Decode(T) error
}

type GobTxEncoder struct {
	W io.Writer
}

func NewGobEncoder(w io.Writer) *GobTxEncoder {
	gob.Register(elliptic.P256())
	gob.Register(&crypto.PublicKey{})
	gob.Register(&crypto.PrivateKey{})
	gob.Register(&crypto.Signature{})
	return &GobTxEncoder{
		w,
	}
}

func (e *GobTxEncoder) Encode(tx *Transaction) error {
	return gob.NewEncoder(e.W).Encode(tx)
}

type GobTxDecoder struct {
	R io.Reader
}

func NewGobDecoder(r io.Reader) *GobTxDecoder {
	gob.Register(elliptic.P256())
	gob.Register(&crypto.PublicKey{})
	gob.Register(&crypto.PrivateKey{})
	gob.Register(&crypto.Signature{})
	return &GobTxDecoder{
		r,
	}
}

func (d *GobTxDecoder) Decode(tx *Transaction) error {
	return gob.NewDecoder(d.R).Decode(tx)
}

/***************************************************************
 * Arquivo: encoding.go
 * Descrição: Implementação de codificação e decodificação.
 * Autor: JoaoRafa19
 * Data de criação: 2024-2025
 * Versão: 0.0.1
 * Licença: MIT License
 * Observações:
 ***************************************************************/

package core

import (
	"crypto/elliptic"
	"encoding/gob"
	"io"
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

func NewGobTxEncoder(w io.Writer) *GobTxEncoder {
	return &GobTxEncoder{
		W: w,
	}
}

func (e *GobTxEncoder) Encode(tx *Transaction) error {
	result := gob.NewEncoder(e.W).Encode(tx)
	return result
}

type GobTxDecoder struct {
	R io.Reader
}

func NewGobDecoder(r io.Reader) *GobTxDecoder {
	return &GobTxDecoder{
		R: r,
	}
}

func (d *GobTxDecoder) Decode(tx *Transaction) error {
	return gob.NewDecoder(d.R).Decode(tx)
}

type GobBlockEncoder struct {
	w io.Writer
}

func NewGobBlockEncoder(w io.Writer) *GobBlockEncoder {
	return &GobBlockEncoder{
		w,
	}
}

func (e *GobBlockEncoder) Encode(block *Block) error {
	return gob.NewEncoder(e.w).Encode(block)
}

type GobBlockDecoder struct {
	r io.Reader
}

func NewGobBlockDecoder(r io.Reader) *GobBlockDecoder {
	return &GobBlockDecoder{
		r,
	}
}

func (d *GobBlockDecoder) Decode(block *Block) error {
	return gob.NewDecoder(d.r).Decode(block)
}

func init() {
	gob.Register(elliptic.P256())
}

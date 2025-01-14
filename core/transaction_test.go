/***************************************************************
 * Arquivo: transaction_test.go
 * Descrição: Testes para a estrutura de transação.
 * Autor: JoaoRafa19
 * Data de criação: 2024-2025
 * Versão: 0.0.1
 * Licença: MIT License
 * Observações: 
 ***************************************************************/

package core

import (
	"bytes"
	"testing"

	"github.com/JoaoRafa19/crypto-go/crypto"
	"github.com/stretchr/testify/assert"
)

func TestSignTransaction(t *testing.T) {
	privateKey := crypto.GeneratePrivateKey()

	tx := &Transaction{
		Data: []byte("foo bar baz"),
	}

	assert.Nil(t, tx.Sign(privateKey))
	assert.NotNil(t, tx.Signature)

}

func TestVerifyTransaction(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: []byte("foo bar baz"),
	}

	assert.Nil(t, tx.Sign(privKey))
	assert.Nil(t, tx.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	tx.From = otherPrivKey.PublicKey()

	assert.NotNil(t, tx.Verify())
}

func TestTxEncodeDecode(t *testing.T) {
	tx := randomTxWithSignature(t)
	buf := &bytes.Buffer{}
	assert.Nil(t, tx.Encode(NewGobEncoder(buf)))
	
	txDecoded := new(Transaction)
	assert.Nil(t, txDecoded.Decode(NewGobDecoder(buf)))
	assert.Equal(t, tx, txDecoded)
}
func randomTxWithSignature(t *testing.T) *Transaction {
	privKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: []byte("foo bar baz"),
	}

	assert.Nil(t, tx.Sign(privKey))
	return tx
}

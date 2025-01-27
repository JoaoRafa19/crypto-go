package core

import (
	"testing"
	"time"

	"github.com/JoaoRafa19/crypto-go/crypto"
	"github.com/JoaoRafa19/crypto-go/types"
	"github.com/stretchr/testify/assert"
)

func TestSignBlock(t *testing.T) {
	priv := crypto.GeneratePrivateKey()
	b := randomBlock(t, 0, types.Hash{})

	assert.Nil(t, b.Sign(priv))
	assert.NotNil(t, b.Signature)
}

func TestVerifyBlock(t *testing.T) {
	priv := crypto.GeneratePrivateKey()
	b := randomBlock(t, 0, types.Hash{})

	assert.Nil(t, b.Sign(priv))
	assert.Nil(t, b.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	b.Validator = otherPrivKey.PublicKey()

	assert.NotNil(t, b.Verify())
	b.Height = 100
	assert.NotNil(t, b.Verify())
}

func randomBlock(t *testing.T, height uint32, prevBlockHas types.Hash) *Block {
	privKey := crypto.GeneratePrivateKey()
	tx := randomTxWithSignature(t)
	header := &Header{
		Version:       1,
		PrevBlockHash: prevBlockHas,
		Height:        height,
		Timestamp:     uint64(time.Now().UnixNano()),
	}

	b, err := NewBlock(header, []Transaction{tx})
	assert.Nil(t, err)

	dataHash, err := CalculateDataHash(b.Transactions)

	assert.Nil(t, err)
	b.Header.DataHash = dataHash
	assert.Nil(t, b.Sign(privKey))

	return b
}

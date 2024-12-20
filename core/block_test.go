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
	b := randomBlock(0, types.Hash{})

	assert.Nil(t, b.Sign(priv))
	assert.NotNil(t, b.Signature)
}

func TestVerifyBlock(t *testing.T) {
	priv := crypto.GeneratePrivateKey()
	b := randomBlock(0, types.Hash{})

	assert.Nil(t, b.Sign(priv))
	assert.Nil(t, b.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	b.Validator = otherPrivKey.PublicKey()

	assert.NotNil(t, b.Verify())
	b.Height = 100
	assert.NotNil(t, b.Verify())
}

func randomBlockWithSignature(t *testing.T, heigh uint32, prevBlockHas types.Hash) *Block {
	privkey := crypto.GeneratePrivateKey()

	bc := randomBlock(heigh, prevBlockHas)
	tx := randomTxWithSignature(t)
	bc.AddTransaction(tx)

	assert.Nil(t, bc.Sign(privkey))
	return bc
}
func randomBlock(height uint32, prevBlockHas types.Hash) *Block {
	header := &Header{
		Version:       1,
		PrevBlockHash: prevBlockHas,
		Height:        height,
		Timestamp:     uint64(time.Now().UnixNano()),
	}

	return NewBlock(header, []Transaction{})
}

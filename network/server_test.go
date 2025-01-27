package network

import (
	"testing"

	"github.com/JoaoRafa19/crypto-go/core"
	"github.com/JoaoRafa19/crypto-go/crypto"
	"github.com/stretchr/testify/assert"
)

func TestProcessMessage_ValidTransaction(t *testing.T) {
	server := NewServer(ServerOpts{})
	privKey := crypto.GeneratePrivateKey()
	tx := core.NewTransaction(nil)
	tx.Sign(privKey)

	decodedMsg := &DecodedMessage{
		From: "testAddr",
		Data: tx,
	}

	err := server.ProcessMessage(decodedMsg)
	assert.Nil(t, err)
	assert.Equal(t, 1, server.MemPool.Len())
	assert.True(t, server.MemPool.Contains(tx.Hash(core.TxHasher{})))
}

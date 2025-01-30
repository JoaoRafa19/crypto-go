package network

import (
	"testing"

	"github.com/JoaoRafa19/crypto-go/core"
	"github.com/JoaoRafa19/crypto-go/crypto"
	"github.com/stretchr/testify/assert"
)

func TestProcessMessage_ValidTransaction(t *testing.T) {
	server, err := NewServer(ServerOpts{})
	assert.Nil(t, err)
	privKey := crypto.GeneratePrivateKey()
	tx := core.NewTransaction(nil)
	tx.Sign(privKey)

	decodedMsg := &DecodedMessage{
		From: "testAddr",
		Data: tx,
	}

	err = server.ProcessMessage(decodedMsg)
	assert.Nil(t, err)
	// assert.Equal(t, 1, server.mempool.Len())
	assert.True(t, server.mempool.Contains(tx.Hash(core.TxHasher{})))
}

package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_KeyPairSignVerifySuccess(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.PublicKey()

	msg := []byte("Hello World!")

	sig, err := privKey.Sign(msg)

	assert.Nil(t, err)
	assert.True(t, sig.Verify(pubKey, msg))

}
func Test_KeyPairSignVerifyFail(t *testing.T) {
	privKey := GeneratePrivateKey()

	msg := []byte("Hello World!")

	sig, err := privKey.Sign(msg)
	assert.Nil(t, err) 	

	otherKey := GeneratePrivateKey()
	otherPub := otherKey.PublicKey()
	otherMessage := []byte("xxxx")

	assert.False(t, sig.Verify(otherPub, msg))

	assert.False(t, sig.Verify(otherPub, otherMessage))

}

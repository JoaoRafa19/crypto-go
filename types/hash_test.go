package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorHashFromBytes(t *testing.T) {
	smalAddress := []byte("test")
	assert.Panics(t, func() { HashFromBytes(smalAddress) })

	bigAddress := []byte("12345678901234567890123123123123123")
	assert.Panics(t, func() { HashFromBytes(bigAddress) })

}

func TestReturnHash(t *testing.T) {
	message := []byte("12345678901234567890123123123123")
	hash := HashFromBytes(message)
	assert.IsType(t, Hash{}, hash)
	assert.False(t, hash.IsZero())
	assert.Len(t, hash, 32)
}
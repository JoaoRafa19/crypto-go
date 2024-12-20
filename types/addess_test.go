package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorAssetsFromBytes(t *testing.T) {
	smalAddress := []byte("test")
	assert.Panics(t, func() { AddressFromBytes(smalAddress) })

	bigAddress := []byte("12345678901234567890123")
	assert.Panics(t, func() { AddressFromBytes(bigAddress) })

}

func TestReturnAaddress(t *testing.T) {
	message := []byte("12345678901234567890")
	addr := AddressFromBytes(message)
	assert.IsType(t, Address{}, addr)
}
func TestReturnCorrectAddress(t *testing.T) {
	message := []byte("addrestest20caractes")
	addr := AddressFromBytes(message)
	assert.IsType(t, Address{}, addr)
	assert.Equal(t, message, addr.ToSlice())
}

package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/JoaoRafa19/crypto-go/types"
)

type PrivateKey struct {
	Key *ecdsa.PrivateKey
}

func (k PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.Key, data)
	if err != nil {
		return nil, err
	}

	return &Signature{
		S: s,
		R: r,
	}, nil
}

func GeneratePrivateKey() PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return PrivateKey{
		Key: key,
	}
}

type PublicKey struct {
	Key *ecdsa.PublicKey
}

func (k PrivateKey) PublicKey() PublicKey {
	return PublicKey{
		Key: &k.Key.PublicKey,
	}
}

func (k PrivateKey) ToBytes() []byte {
	return k.Key.D.Bytes()
}

func PrivateKeyFromBytes(data []byte) (PrivateKey, error) {
	d := new(big.Int).SetBytes(data)
	privKey := new(ecdsa.PrivateKey)
	privKey.PublicKey.Curve = elliptic.P256()
	privKey.D = d
	privKey.PublicKey.X, privKey.PublicKey.Y = elliptic.P256().ScalarBaseMult(d.Bytes())
	return PrivateKey{Key: privKey}, nil
}

func (k PublicKey) ToBytes() []byte {
	return elliptic.MarshalCompressed(k.Key, k.Key.X, k.Key.Y)
}

func PublicKeyFromBytes(data []byte) (PublicKey, error) {
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), data)
	if x == nil {
		return PublicKey{}, fmt.Errorf("invalid public key data")
	}
	return PublicKey{Key: &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}}, nil
}

func (k PublicKey) ToSlice() []byte {
	return elliptic.MarshalCompressed(k.Key, k.Key.X, k.Key.Y)

}

func (k PublicKey) Address() types.Address {
	h := sha256.Sum256(k.ToSlice())

	return types.AddressFromBytes(h[len(h)-20:])
}

type Signature struct {
	S, R *big.Int
}

func (sig Signature) Verify(pubKey PublicKey, data []byte) bool {
	return ecdsa.Verify(pubKey.Key, data, sig.R, sig.S)
}

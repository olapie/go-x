package xsecurity

import (
	"crypto/ecdsa"
	"crypto/rand"
	"testing"
	"time"

	"go.olapie.com/x/xtest"
)

func TestEncodePrivateKey(t *testing.T) {
	ecdsaKey, err := GeneratePrivateKey(EcdsaP256)
	xtest.NoError(t, err)
	pk := ecdsaKey.(*ecdsa.PrivateKey)
	data, err := EncodePrivateKey(pk, "hello")
	xtest.NoError(t, err)
	_, err = DecodePrivateKey(data, "hi")
	xtest.Error(t, err)
	ecdsaKey2, err := DecodePrivateKey(data, "hello")
	xtest.NoError(t, err)
	pk2 := ecdsaKey2.(*ecdsa.PrivateKey)
	digest := []byte(SHA1(time.Now().String()))
	sign1, err := ecdsa.SignASN1(rand.Reader, pk, digest[:])
	xtest.NoError(t, err)
	sign2, err := ecdsa.SignASN1(rand.Reader, pk2, digest[:])
	xtest.NoError(t, err)

	xtest.True(t, ecdsa.VerifyASN1(&pk.PublicKey, digest[:], sign1))
	xtest.True(t, ecdsa.VerifyASN1(&pk.PublicKey, digest[:], sign2))
	xtest.True(t, ecdsa.VerifyASN1(&pk2.PublicKey, digest[:], sign1))
	xtest.True(t, ecdsa.VerifyASN1(&pk2.PublicKey, digest[:], sign2))

	xtest.Equal(t, pk.X, pk2.X)
	xtest.Equal(t, pk.Y, pk2.Y)
	xtest.Equal(t, pk.D, pk2.D)
}

func TestEncodePublicKey(t *testing.T) {
	ecdsaKey, err := GeneratePrivateKey(EcdsaP256)
	xtest.NoError(t, err)
	pk := ecdsaKey.(*ecdsa.PrivateKey)
	data, err := EncodePublicKey(&pk.PublicKey)
	xtest.NoError(t, err)
	pub, err := DecodePublicKey(data)
	xtest.NoError(t, err)
	digest := []byte(SHA1(time.Now().String()))
	sign, err := ecdsa.SignASN1(rand.Reader, pk, digest[:])
	xtest.NoError(t, err)
	xtest.True(t, ecdsa.VerifyASN1(pub.(*ecdsa.PublicKey), digest[:], sign))
	xtest.True(t, pk.PublicKey.Equal(pub))
}

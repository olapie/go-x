package xjwt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.olapie.com/x/xtest"
)

func TestValidToken(t *testing.T) {
	pri, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	pub := pri.Public().(*ecdsa.PublicKey)
	userID := uuid.NewString()[:8]
	appID := uuid.NewString()[:8]
	expiresAt := time.Now().Add(time.Second * 2)
	tokenString, err := Sign(pri, appID, userID, expiresAt)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tokenString)
	parsedAppID, parsedUserID, err := Parse(pub, tokenString)
	if err != nil {
		t.Fatal(err)
	}
	xtest.Equal(t, appID, parsedAppID)
	xtest.Equal(t, userID, parsedUserID)

	time.Sleep(time.Second * 2)
	_, _, err = Parse(pub, tokenString)
	if err == nil {
		t.Fatal("should be expired")
	} else {
		t.Log(err)
	}
}

func TestExpiredToken(t *testing.T) {
	pri, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	pub := pri.Public().(*ecdsa.PublicKey)
	userID := uuid.NewString()[:8]
	appID := uuid.NewString()[:8]
	expiresAt := time.Now().Add(time.Second * 2)
	tokenString, err := Sign(pri, appID, userID, expiresAt)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 2)
	_, _, err = Parse(pub, tokenString)
	if err == nil {
		t.Fatal("should be expired")
	} else {
		t.Log(err)
	}
}

func TestInvalidKey(t *testing.T) {
	pri, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	pri2, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	pub := pri2.Public().(*ecdsa.PublicKey)
	userID := uuid.NewString()[:8]
	appID := uuid.NewString()[:8]
	expiresAt := time.Now().Add(time.Second * 2)
	tokenString, err := Sign(pri, appID, userID, expiresAt)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = Parse(pub, tokenString)
	if err == nil {
		t.Fatal("should be invalid key")
	} else {
		t.Log(err)
	}
}

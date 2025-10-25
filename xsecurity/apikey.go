package xsecurity

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// GenerateAPIKey generates API key for each request
// privateKey is pre-configured in each officially built client package
func GenerateAPIKey(privateKey *ecdsa.PrivateKey, traceID string, timestamp int64) (string, error) {
	digest := generateRequestDigest(traceID, timestamp)
	signedDigest, err := ecdsa.SignASN1(rand.Reader, privateKey, digest)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(signedDigest), nil
}

func VerifyAPIKey(publicKey *ecdsa.PublicKey, traceID string, timestamp int64, apiKey string) (bool, error) {
	signedDigest, err := base64.RawURLEncoding.DecodeString(apiKey)
	if err != nil {
		return false, fmt.Errorf("base64 url decode: %w", err)
	}
	digest := generateRequestDigest(traceID, timestamp)
	return ecdsa.VerifyASN1(publicKey, digest, signedDigest), nil
}

func generateRequestDigest(traceID string, timestamp int64) []byte {
	msg := fmt.Sprintf("%s.%d", traceID, timestamp)
	sum := sha256.Sum256([]byte(msg))
	return sum[:]
}

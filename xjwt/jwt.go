package xjwt

import (
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func Sign(privateKey *ecdsa.PrivateKey, appID, userID string, expiresAt time.Time) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    appID,
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}
	//SigningMethodES256: ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	//SigningMethodES512: ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	//Long key generates long token
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(privateKey)
}

func Parse(publicKey *ecdsa.PublicKey, tokenString string) (appID, userID string, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	}) // e.g. option jwt.WithLeeway(5*time.Second)
	if err != nil {
		return "", "", fmt.Errorf("parse token: %s, %w", tokenString, err)
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims.Issuer, claims.Subject, nil
	}
	return "", "", fmt.Errorf("invalid token: %s", tokenString)
}

package xjwt

import (
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.olapie.com/x/xerror"
	"go.olapie.com/x/xhttpheader"
	"go.olapie.com/x/xtype"
)

const (
	ErrNoAuthorization = xerror.String("no authorization token")
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

func VerifyAuthorization[T ~map[string][]string | ~map[string]string](pubKey *ecdsa.PublicKey, h T) (*xtype.AuthResult, error) {
	accessToken := xhttpheader.GetBearer(h)
	if accessToken == "" {
		accessToken = xhttpheader.GetAuthorization(h)
	}

	if accessToken == "" {
		return nil, ErrNoAuthorization
	}

	appID, uid, err := Parse(pubKey, accessToken)
	if err != nil {
		return nil, xerror.Unauthorized("failed to parse access token: %s", err)
	}

	return &xtype.AuthResult{
		AppID:  appID,
		UserID: xtype.NewUserID(uid),
	}, nil
}

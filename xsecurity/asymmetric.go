package xsecurity

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"time"
)

type AsymmetricAlgorithm int

const (
	_ AsymmetricAlgorithm = iota
	RSA2048
	RSA4096
	Ed25519
	EcdsaP224
	EcdsaP256
	EcdsaP384
	EcdsaP521
)

func (a AsymmetricAlgorithm) String() string {
	switch a {
	case RSA2048:
		return "RSA2048"
	case RSA4096:
		return "RSA4096"
	case Ed25519:
		return "Ed25519"
	case EcdsaP224:
		return "EcdsaP224"
	case EcdsaP256:
		return "EcdsaP256"
	case EcdsaP384:
		return "EcdsaP384"
	case EcdsaP521:
		return "EcdsaP521"
	default:
		return fmt.Sprint(int(a))
	}
}

type PrivateKeySet interface {
	*rsa.PrivateKey | *ed25519.PrivateKey | *ecdsa.PrivateKey
}

type PublicKeyTypeSet interface {
	*rsa.PublicKey | *ed25519.PublicKey | *ecdsa.PublicKey
}

// GeneratePrivateKey generates private key with the type
// Return type would be *rsa.PrivateKey, *ed25519.PrivateKey or *ecdsa.PrivateKey
func GeneratePrivateKey(a AsymmetricAlgorithm) (any, error) {
	switch a {
	case RSA2048:
		return rsa.GenerateKey(rand.Reader, 2048)
	case RSA4096:
		return rsa.GenerateKey(rand.Reader, 4096)
	case Ed25519:
		_, priv, err := ed25519.GenerateKey(rand.Reader)
		return priv, err
	case EcdsaP224:
		return ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case EcdsaP256:
		return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case EcdsaP384:
		return ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case EcdsaP521:
		return ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		return nil, fmt.Errorf("invalid private key type: %d", int(a))
	}
}

// GetPublicKey get public key of private key of types: *rsa.PrivateKey, *ed25519.PrivateKey or *ecdsa.PrivateKey
func GetPublicKey(priv any) any {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	case ed25519.PrivateKey:
		return k.Public().(ed25519.PublicKey)
	default:
		return nil
	}
}

func EncodePublicKey(pub any) ([]byte, error) {
	return x509.MarshalPKIXPublicKey(pub)
}

func EncodePrivateKey(priv any, passphrase string) ([]byte, error) {
	data, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, err
	}
	return Encrypt(data, passphrase)
}

func DecodePrivateKey(data []byte, passphrase string) (any, error) {
	size := len(data) - HeaderSize
	if size < 0 {
		return nil, fmt.Errorf("invalid data")
	}
	raw, err := Decrypt(data, passphrase)
	if err != nil {
		return nil, err
	}
	return x509.ParsePKCS8PrivateKey(raw)
}

func DecodePublicKey(data []byte) (any, error) {
	return x509.ParsePKIXPublicKey(data)
}

type CertificateOptions struct {
	Algorithm    AsymmetricAlgorithm
	Hosts        []string
	IsCA         bool
	NotBefore    time.Time
	NotAfter     time.Time
	Organization string
}

func GenerateCertificate(optFns ...func(options *CertificateOptions)) (privateKey, publicKey []byte, err error) {
	options := &CertificateOptions{
		Algorithm:    EcdsaP256,
		Hosts:        []string{"127.0.0.1"},
		IsCA:         false,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(100, 0, 0),
		Organization: "olapie",
	}

	for _, fn := range optFns {
		fn(options)
	}

	priv, err := GeneratePrivateKey(options.Algorithm)
	if err != nil {
		return nil, nil, err
	}

	keyUsage := x509.KeyUsageDigitalSignature
	if _, isRSA := priv.(*rsa.PrivateKey); isRSA {
		keyUsage |= x509.KeyUsageKeyEncipherment
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate serial number: %v", err)
	}

	certTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{options.Organization},
		},
		NotBefore:             options.NotBefore,
		NotAfter:              options.NotAfter,
		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	for _, host := range options.Hosts {
		if ip := net.ParseIP(host); ip != nil {
			certTemplate.IPAddresses = append(certTemplate.IPAddresses, ip)
		} else {
			certTemplate.DNSNames = append(certTemplate.DNSNames, host)
		}
	}

	if options.IsCA {
		certTemplate.IsCA = true
		certTemplate.KeyUsage |= x509.KeyUsageCertSign
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &certTemplate, &certTemplate, GetPublicKey(priv), priv)
	if err != nil {
		err = fmt.Errorf("Failed to create certificate: %v", err)
		return
	}

	var pubKeyBuf bytes.Buffer
	err = pem.Encode(&pubKeyBuf, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		err = fmt.Errorf("pem.Encode public key: %v", err)
		return
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		err = fmt.Errorf("x509.MarshalPKCS8PrivateKey: %v", err)
		return
	}
	var privKeyBuf bytes.Buffer
	err = pem.Encode(&privKeyBuf, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	if err != nil {
		err = fmt.Errorf("pem.Encode private key: %v", err)
		return
	}
	return pubKeyBuf.Bytes(), privKeyBuf.Bytes(), nil
}

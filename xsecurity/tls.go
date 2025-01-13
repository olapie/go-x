package xsecurity

import (
	"crypto/tls"
	"fmt"
)

func NewTLSConfig(optFns ...func(options *CertificateOptions)) (*tls.Config, error) {
	pubKey, privKey, err := GenerateCertificate(optFns...)
	if err != nil {
		return nil, fmt.Errorf("generate certificate: %w", err)
	}

	cert, err := tls.X509KeyPair(pubKey, privKey)
	if err != nil {
		return nil, fmt.Errorf("tls.X509KeyPair: %w", err)
	}

	return &tls.Config{Certificates: []tls.Certificate{cert}}, nil
}

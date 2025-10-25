package xsecurity

import "os"

// SecretSeed is set by linker
var SecretSeed string

func GetSecretSeed() string {
	if SecretSeed != "" {
		return SecretSeed
	}
	return os.Getenv("OLA_SECRET_SEED")
}

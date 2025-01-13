package xrand

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"

	"github.com/google/uuid"
)

var defaultOTPGenerator *otpGenerator

const digitAlphaString = "23456789ABCDEFGHIJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz"

type otpGenerator struct {
	mu       sync.Mutex
	lenToMax map[int]*big.Int
}

func newOTPGenerator() *otpGenerator {
	return &otpGenerator{
		lenToMax: map[int]*big.Int{},
	}
}

func (s *otpGenerator) generateDigitCode(length int) (string, error) {
	maxVal := s.getMaxRand(length)
	var code string
	for len(code) < length {
		n, err := rand.Int(rand.Reader, maxVal)
		if err != nil {
			return "", fmt.Errorf("rand.Int: %w", err)
		}
		code += n.String()
	}
	return code[:length], nil
}

func (s *otpGenerator) generateAlphaDigitCode(length int) (string, error) {
	maxVal := big.NewInt(int64(len(digitAlphaString)))
	a := make([]byte, length)
	for i := range a {
		index, err := rand.Int(rand.Reader, maxVal)
		if err != nil {
			return "", fmt.Errorf("rand.Int: %w", err)
		}
		a[i] = digitAlphaString[index.Int64()]
	}
	return string(a), nil
}

func (s *otpGenerator) getMaxRand(length int) *big.Int {
	maxVal := s.lenToMax[length]
	if maxVal != nil {
		return maxVal
	}
	s.mu.Lock()
	maxVal = s.lenToMax[length]
	if maxVal != nil {
		return maxVal
	}
	maxVal = big.NewInt(1)
	n10 := big.NewInt(10)
	// length=6 => max=1e7
	// length=8 => max=1e9
	for i := 0; i < length; i++ {
		maxVal.Mul(maxVal, n10)
	}
	s.lenToMax[length] = maxVal
	s.mu.Unlock()
	return maxVal
}

func DigitCode(n int) string {
	if defaultOTPGenerator == nil {
		defaultOTPGenerator = newOTPGenerator()
	}
	v, err := defaultOTPGenerator.generateDigitCode(n)
	if err != nil {
		panic(err)
	}
	return v
}

func AlphaDigitCode(n int) string {
	if defaultOTPGenerator == nil {
		defaultOTPGenerator = newOTPGenerator()
	}
	v, err := defaultOTPGenerator.generateAlphaDigitCode(n)
	if err != nil {
		panic(err)
	}
	return v
}

func Bytes(n int) []byte {
	if n <= 0 {
		return []byte{}
	}

	b := make([]byte, n)
	_, err := rand.Read(b)
	if err == nil {
		return b
	}

	p := b
	for len(p) > 0 {
		id := uuid.New()
		n := copy(p, id[:])
		p = p[n:]
	}
	return b
}

func String(n int) string {
	return StringT[string](n)
}

func StringT[T ~string](n int) T {
	if n <= 0 {
		return ""
	}

	b := Bytes(n + 1/2)
	s := fmt.Sprintf("%x", b)
	return T(s[:n])
}

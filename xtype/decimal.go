package xtype

import (
	"encoding"
	"fmt"
	"math/big"
	"strings"
)

var _ encoding.TextMarshaler = (*Decimal)(nil)
var _ encoding.TextMarshaler = Decimal{}
var _ encoding.TextUnmarshaler = (*Decimal)(nil)

type Decimal struct {
	Int *big.Int
	Exp int32
}

func NewDecimalFromString(s string) (*Decimal, error) {
	if s == "" || s == "." {
		return nil, fmt.Errorf("invalid string")
	}

	if s == "0" {
		return &Decimal{
			Int: big.NewInt(0),
			Exp: 0,
		}, nil
	}

	n := len(s)

	pointPos := -1
	for i := range n {
		if s[i] == '.' {
			if pointPos >= 0 {
				return nil, fmt.Errorf("unexpected decimal point at %d", i)
			}
			pointPos = i
		} else {
			if s[i] < '0' || s[i] > '9' {
				return nil, fmt.Errorf("invalid digit %v at %d", s[i], i)
			}
		}
	}

	if pointPos < 0 {
		intVal, ok := big.NewInt(0).SetString(s, 10)
		if !ok {
			return nil, fmt.Errorf("invalid integer part")
		}
		return &Decimal{
			Int: intVal,
			Exp: 0,
		}, nil
	}

	integerPart := trimLeadingZeros(s[:pointPos])
	fractionalPart := trimTrailingZeros(s[pointPos+1:])

	v, ok := big.NewInt(0).SetString(integerPart+fractionalPart, 10)
	if !ok {
		return nil, fmt.Errorf("invalid decimal")
	}

	return &Decimal{
		Int: v,
		Exp: int32(-len(fractionalPart)),
	}, nil
}

func (d *Decimal) String() string {
	s := d.Int.String()
	if d.Exp == 0 {
		return s
	}

	if d.Exp > 0 {
		return s + strings.Repeat("0", int(d.Exp))
	}
	k := len(s) + int(d.Exp)
	if k <= 0 {
		return fmt.Sprintf("0.%s%s", strings.Repeat("0", -k), s)
	}
	return fmt.Sprintf("%s.%s", s[:k], s[k:])
}

func (d Decimal) MarshalText() (text []byte, err error) {
	return []byte(d.String()), nil
}

func (d *Decimal) UnmarshalText(text []byte) error {
	v, err := NewDecimalFromString(string(text))
	if err != nil {
		return err
	}
	d.Int = v.Int
	d.Exp = v.Exp
	return nil
}

func trimLeadingZeros(s string) string {
	for i, c := range s {
		if c != '0' {
			return s[i:]
		}
	}
	return s
}

func trimTrailingZeros(s string) string {
	i := len(s) - 1
	for i >= 0 && s[i] == '0' {
		i--
	}
	return s[0 : i+1]
}

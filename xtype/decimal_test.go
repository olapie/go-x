package xtype

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
)

func TestNewDecimalFromString(t *testing.T) {
	type TestCase struct {
		String        string
		DecimalString string
	}

	happyCases := []TestCase{
		{
			String:        "0",
			DecimalString: "0",
		},
		{
			String:        "0000000",
			DecimalString: "0",
		},
		{
			String:        "0.0",
			DecimalString: "0",
		},
		{
			String:        "0.00000",
			DecimalString: "0",
		},
		{
			String:        "000.00000",
			DecimalString: "0",
		},
		{
			String:        "000.0",
			DecimalString: "0",
		},
		{
			String:        "000.",
			DecimalString: "0",
		},
		{
			String:        "0.1",
			DecimalString: "0.1",
		},
		{
			String:        "0.0001",
			DecimalString: "0.0001",
		},
		{
			String:        "0.0001000",
			DecimalString: "0.0001",
		},
		{
			String:        "1.0001000",
			DecimalString: "1.0001",
		},
		{
			String:        ".0001000",
			DecimalString: "0.0001",
		},
		{
			String:        "123458023890183018230138",
			DecimalString: "123458023890183018230138",
		},
		{
			String:        "123458023890183018230138.",
			DecimalString: "123458023890183018230138",
		},
		{
			String:        "123458023890183018230138.1983018301831038013",
			DecimalString: "123458023890183018230138.1983018301831038013",
		},
		{
			String:        "123458023890183018230138.000001983018301831038013",
			DecimalString: "123458023890183018230138.000001983018301831038013",
		},
		{
			String:        "00000123458023890183018230138.000001983018301831038013",
			DecimalString: "123458023890183018230138.000001983018301831038013",
		},
	}

	for _, test := range happyCases {
		d, err := NewDecimalFromString(test.String)
		if err != nil {
			t.Fatalf("%s: %v", test.String, err)
		}

		got := d.String()
		if got != test.DecimalString {
			t.Log(test.String, d.Int.String(), d.Exp)
			t.Fatalf("expected %s, got %s", test.DecimalString, got)
		}
	}

	badCases := []string{"", ".", "0.1.", ".1.", "a", "1.a", "1.2.3.0", "00000.11111000.0"}
	for _, test := range badCases {
		d, err := NewDecimalFromString(test)
		if err == nil {
			t.Fatalf("expected error for %s", test)
		}

		if d != nil {
			t.Fatalf("expected nil for %s", test)
		}
	}
}

func TestDecimalJSONMarshal(t *testing.T) {
	type Item struct {
		Price Decimal `json:"price"`
	}

	decimalString := fmt.Sprint(rand.Float64())
	decimal, err := NewDecimalFromString(decimalString)
	if err != nil {
		t.Fatal(err)
	}

	jsonBytes, err := json.Marshal(Item{
		Price: *decimal,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(jsonBytes))

	var item Item
	err = json.Unmarshal(jsonBytes, &item)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(item.Price.String())
	if item.Price.String() != decimalString {
		t.Fatal(item.Price.String(), decimalString)
	}
}

func TestDecimalJSONMarshal_Pointer(t *testing.T) {
	type Item struct {
		Price *Decimal `json:"price"`
	}

	decimalString := fmt.Sprint(rand.Float64())
	decimal, err := NewDecimalFromString(decimalString)
	if err != nil {
		t.Fatal(err)
	}

	jsonBytes, err := json.Marshal(Item{
		Price: decimal,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(jsonBytes))

	var item Item
	err = json.Unmarshal(jsonBytes, &item)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(item.Price.String())
	if item.Price.String() != decimalString {
		t.Fatal(item.Price.String(), decimalString)
	}
}

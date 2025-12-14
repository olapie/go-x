package naming

import "testing"

func TestToPascalWithAcronym(t *testing.T) {
	type TestCase struct {
		Input  string
		Output string
	}
	tests := []TestCase{
		{"-----", ""},
		{"---._ \t", ""},
		{"--b-._ \ta", "BA"},
		{"b--b-._ \ta", "BBA"},
		{"b--bC-._ \ta", "BBCA"},
		{"bef--bC-._ \ta", "BefBCA"},
		{"xml___http", "XMLHTTP"},
		{"xml.http", "XMLHTTP"},
		{"_xml.http", "XMLHTTP"},
		{"__xml.http__", "XMLHTTP"},
	}
	for _, test := range tests {
		output := ToPascal(test.Input, WithAcronym())
		if output != test.Output {
			t.Fatalf("Test %s, got %s, want %s", test.Input, output, test.Output)
		}
	}
}

func TestToPascalWithoutAcronym(t *testing.T) {
	type TestCase struct {
		Input  string
		Output string
	}
	tests := []TestCase{
		{"-----", ""},
		{"---._ \t", ""},
		{"--b-._ \ta", "BA"},
		{"b--b-._ \ta", "BBA"},
		{"b--bC-._ \ta", "BBCA"},
		{"bef--bC-._ \ta", "BefBCA"},
		{"xml___http", "XmlHttp"},
		{"xml.http", "XmlHttp"},
		{"_xml.http", "XmlHttp"},
		{"__xml.http__", "XmlHttp"},
	}
	for _, test := range tests {
		output := ToPascal(test.Input)
		if output != test.Output {
			t.Fatalf("Test %s, got %s, want %s", test.Input, output, test.Output)
		}
	}
}

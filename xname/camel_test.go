package naming

import "testing"

func TestToCamelWithAcronym(t *testing.T) {
	type TestCase struct {
		Input  string
		Output string
	}
	tests := []TestCase{
		{"-----", ""},
		{"---._ \t", ""},
		{"--b-._ \ta", "bA"},
		{"b--b-._ \ta", "bBA"},
		{"b--bC-._ \ta", "bBCA"},
		{"bef--bC-._ \ta", "befBCA"},
		{"xml___http", "xmlHTTP"},
		{"xml.http", "xmlHTTP"},
		{"_xml.http", "xmlHTTP"},
		{"__xml.http__", "xmlHTTP"},
		{"URLString", "urlString"},
	}
	for _, test := range tests {
		output := ToCamel(test.Input, WithAcronym())
		if output != test.Output {
			t.Fatalf("Test %s, got %s, want %s", test.Input, output, test.Output)
		}
	}
}

func TestToCamelWithoutAcronym(t *testing.T) {
	type TestCase struct {
		Input  string
		Output string
	}
	tests := []TestCase{
		{"-----", ""},
		{"---._ \t", ""},
		{"--b-._ \ta", "bA"},
		{"b--b-._ \ta", "bBA"},
		{"b--bC-._ \ta", "bBCA"},
		{"bef--bC-._ \ta", "befBCA"},
		{"xml___http", "xmlHttp"},
		{"xml.http", "xmlHttp"},
		{"_xml.http", "xmlHttp"},
		{"__xml.http__", "xmlHttp"},
		{"URLString", "urlstring"},
		{"doc_url", "docUrl"},
	}
	for _, test := range tests {
		output := ToCamel(test.Input)
		if output != test.Output {
			t.Fatalf("Test %s, got %s, want %s", test.Input, output, test.Output)
		}
	}
}

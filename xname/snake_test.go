package naming

import "testing"

func TestToSnake(t *testing.T) {
	type TestCase struct {
		Input  string
		Output string
	}
	tests := []TestCase{
		{"-----", ""},
		{"---._ \t", ""},
		{"--b-._ \ta", "b_a"},
		{"b--b-._ \ta", "b_b_a"},
		{"b--bC-._ \ta", "b_b_c_a"},
		{"bef--bC-._ \ta", "bef_b_c_a"},
		{"xml___http", "xml_http"},
		{"xml.http", "xml_http"},
		{"_xml.http", "xml_http"},
		{"__xml.http__", "xml_http"},
		{"URLString", "url_string"},
		{"HTTPRequestHandler", "http_request_handler"},
	}
	for _, test := range tests {
		output := ToSnake(test.Input, WithAcronym())
		if output != test.Output {
			t.Fatalf("Test %s, got %s, want %s", test.Input, output, test.Output)
		}
	}
}

func TestToSnakeWithoutAcronym(t *testing.T) {
	type TestCase struct {
		Input  string
		Output string
	}
	tests := []TestCase{
		{"-----", ""},
		{"---._ \t", ""},
		{"--b-._ \ta", "b_a"},
		{"b--b-._ \ta", "b_b_a"},
		{"b--bC-._ \ta", "b_b_c_a"},
		{"bef--bC-._ \ta", "bef_b_c_a"},
		{"xml___http", "xml_http"},
		{"xml.http", "xml_http"},
		{"_xml.http", "xml_http"},
		{"__xml.http__", "xml_http"},
		{"URLString", "urlstring"},
		{"HTTPRequestHandler", "httprequest_handler"},
	}
	for _, test := range tests {
		output := ToSnake(test.Input)
		if output != test.Output {
			t.Fatalf("Test %s, got %s, want %s", test.Input, output, test.Output)
		}
	}
}

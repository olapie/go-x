package naming

import "testing"

func TestNormalize(t *testing.T) {
	type TestCase struct {
		Input  string
		Output string
	}
	tests := []TestCase{
		{"-----", ""},
		{"---._ \t", ""},
		{"--b-._ \ta", "b_a"},
		{"b--b-._ \ta", "b_b_a"},
		{"b--bC-._ \ta", "b_bC_a"},
		{"bef--bC-._ \ta", "bef_bC_a"},
		{"xml___http", "xml_http"},
		{"xml.http", "xml_http"},
		{"_xml.http", "xml_http"},
		{"__xml.http__", "xml_http"},
	}
	for _, test := range tests {
		output := string(normalize(test.Input))
		if output != test.Output {
			t.Fatalf("Test %s, got %s, want %s", test.Input, output, test.Output)
		}
	}
}

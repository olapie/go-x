package xhtml

import (
	"fmt"
	"testing"
)

func TestParseXPath(t *testing.T) {
	segments, err := ParseXPath("/html/body/div[3]/div[2]/div/div/div[6]/div[2]/table/tbody/tr[5]/td[6]")
	if err != nil {
		t.Fatal(err)
	}

	expected := "[{html 1} {body 1} {div 3} {div 2} {div 1} {div 1} {div 6} {div 2} {table 1} {tbody 1} {tr 5} {td 6}]"
	got := fmt.Sprint(segments)
	if expected != got {
		t.Fatalf("expected:\t%s\ngot:\t%s\n", expected, got)
	}
}

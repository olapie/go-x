package xhtml

import (
	"fmt"
	"golang.org/x/net/html"
	"strings"
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

func TestFindByXPath(t *testing.T) {
	htmlString := `
<html>
	<head>
	</head>
	<body>
		<div>
			<div>
			</div>
			<div>
			</div>
			<div>
			</div>
		</div>
		<div id="unused"></div>
		<div id="unused"></div>
		<div>
			<div>
			</div>
			<div id="unused"></div>
			<div>
				<p>I am here</p>
			</div>
			<div>
			</div>
		</div>
	</body>
</html>
`
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		t.Fatal(err)
	}

	node, err := FindByXPath(doc, "/html/body/div[2]/div[2]/p", func(n *html.Node) bool {
		if len(n.Attr) != 0 {
			if n.Attr[0].Key == "id" && n.Attr[0].Val == "unused" {
				return false
			}
		}
		return true
	})

	if err != nil {
		t.Fatal(err)
	}
	if node.Data != "p" {
		t.Fatalf("expect: p, got: %s", node.Data)
	}
	if node.FirstChild == nil {
		t.Fatal("FirstChild is not nil")
	}
	if node.FirstChild.Data != "I am here" {
		t.Fatalf("expect: I am here, got: %s", node.FirstChild.Data)
	}
}

package xhtml

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

var xpathSegmentRegexp = regexp.MustCompile(`(/\w+(\[\d+\])?)?`)

type XPathSegment struct {
	Name     string
	Position int
}

// ParseXPath parse path into XPathSegment
// only support simple path like /html/body/div[3]/div[2]/div/div/div[6]/div[2]/table/tbody/tr[5]/td[6]
func ParseXPath(path string) ([]XPathSegment, error) {
	if path == "" {
		return nil, errors.New("xpath cannot be empty")
	}

	if !xpathSegmentRegexp.MatchString(path) {
		return nil, errors.New("invalid xpath")
	}

	segments := strings.Split(path[1:], "/")
	nodes := make([]XPathSegment, 0, len(segments))
	for _, seg := range segments {
		node := XPathSegment{
			Name:     seg,
			Position: 1,
		}
		left := strings.Index(seg, "[")
		if left > 0 {
			i, err := strconv.ParseInt(seg[left+1:len(seg)-1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid xpath: %s", seg)
			}
			if i < 0 {
				return nil, fmt.Errorf("invalid xpath: %s, position cannot be negative", seg)
			}
			node.Name = seg[:left]
			node.Position = int(i)
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

// FindByXPath finds node with XPath 'path'
// root must be DocumentNode
func FindByXPath(root *html.Node, path string, filter func(node *html.Node) bool) (*html.Node, error) {
	if root.Type != html.DocumentNode {
		return nil, errors.New("root is not document node")
	}

	xpathNodes, err := ParseXPath(path)
	if err != nil {
		return nil, fmt.Errorf("parse xpath: %w", err)
	}

	if len(xpathNodes) == 0 {
		return root, nil
	}

	current := root.FirstChild
	for i, xpathNode := range xpathNodes {
		for j := 0; j < xpathNode.Position && current != nil; current = current.NextSibling {
			if filter != nil && !filter(current) {
				continue
			}

			if current.Data == xpathNode.Name {
				j++
				if j == xpathNode.Position {
					break
				}
			}
		}

		if current == nil || current.Data != xpathNode.Name {
			return nil, errors.New(fmt.Sprintf("failed at %v", xpathNodes[:i+1]))
		}
		current = current.FirstChild
	}
	return current.Parent, nil
}

package main

import (
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"os"
	"regexp"
	"strings"
)

var el_type = flag.String("type", "", "node type")
var el_key = flag.String("key", "", "key name")
var el_val = flag.String("val", "", "value")
var el_not = flag.Bool("not", false, "Inverse filter")
var el_cont = flag.Bool("contains", false, "Match val by contains")
var el_regexp = flag.Bool("regexp", false, "Match val by regexp")

func Match(v string, s string) bool {
	if *el_cont {
		return strings.Contains(v, s)
	}
	if *el_regexp {
		return rex.Match([]byte(v))
	}
	return v == s
}

func Filter(n *html.Node) bool {
	if n.Data != *el_type && *el_type != "" {
		return *el_not
	}
	for _, v := range n.Attr {
		if (*el_key == "" || v.Key == *el_key) &&
			(*el_val == "" || Match(v.Val, *el_val)) {
			return !*el_not
		}
	}
	return *el_not
}

func Body(doc *html.Node) []*html.Node {
	var crawler func(*html.Node)
	var nodes []*html.Node
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && Filter(node) {
			if !*el_not {
				nodes = append(nodes, node)
				return
			}
		}
		for child := node.FirstChild; child != nil; {
			next := child.NextSibling
			if *el_not && child.Type == html.ElementNode && !Filter(child) {
				node.RemoveChild(child)
			} else {
				crawler(child)
			}
			child = next
		}
	}
	if *el_not {
		if (!Filter(doc)) {
			return nodes
		}
		nodes = append(nodes, doc)
	}
	crawler(doc)
	return nodes
}

func renderNodes(n []*html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	for _, v := range n {
		html.Render(w, v)
	}
	return html.UnescapeString(buf.String())
}

var rex *regexp.Regexp

func main() {
	flag.Parse()
	if *el_regexp {
		var err error
		rex, err = regexp.Compile(*el_val)
		if err != nil {
			os.Exit(1)
		}
	}
	doc, _ := html.Parse(os.Stdin)
	nodes := Body(doc)
	body := renderNodes(nodes)
	fmt.Println(body)
}

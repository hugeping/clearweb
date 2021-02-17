package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strings"
)

type Link struct {
	Href string `xml:"href,attr"`
}

type Item struct {
	Title string `xml:"title"`
	Link  Link   `xml:"link"`
	Desc  string `xml:"content"`
	num   int
}

type Atom struct {
	XMLName xml.Name `xml:"http://www.w3.org/2005/Atom feed"`
	Items   []Item   `xml:"entry"`
	Title   string   `xml:"title"`
}

func str_esc(l string) string {
	l = strings.Replace(l, "&", "&amp;", -1)
	l = strings.Replace(l, "<", "&lt;", -1)
	l = strings.Replace(l, ">", "&gt;", -1)
	return l
}

func main() {
	rss := Atom{}
	decoder := xml.NewDecoder(os.Stdin)
	err := decoder.Decode(&rss)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("<rss version='2.0'>\n<channel>\n<title>%s</title>\n", rss.Title)
	for _, v := range rss.Items {
		fmt.Printf("<item>\n<title>%s</title>\n<link>%s</link>\n<description>%s</description>\n</item>\n",
			str_esc(v.Title), v.Link.Href, str_esc(v.Desc))
	}
	fmt.Printf("</channel>\n</rss>\n")
}

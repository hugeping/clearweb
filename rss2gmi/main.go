package main

import (
	"encoding/xml"
	"golang.org/x/net/html/charset"
	"golang.org/x/net/html"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
)

type Item struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Desc  string `xml:"description"`
	Full  string `xml:"encoded"`
	num   int
}

type Channel struct {
	Items []Item `xml:"item"`
	Title string `xml:"title"`
	Desc string `xml:"description"`
}

type Rss struct {
	Channel Channel `xml:"channel"`
}

func html_decode(s string) string {
	TxtContent := ""
	tok := html.NewTokenizer(strings.NewReader(s))
	prev := tok.Token()

	for {
		tt := tok.Next()
		n, _ := tok.TagName()
		nam := string(n)
		switch {
		case tt == html.ErrorToken:
			return TxtContent
		case tt == html.StartTagToken:
			prev = tok.Token()
		case tt == html.TextToken:
			if prev.Data == "script" {
				continue
			}
			txt := strings.Replace(html.UnescapeString(string(tok.Text())), "\r", "", -1)
			txt = strings.Replace(txt, "\n", "", -1)
			TxtContent += txt
		}
		if tt == html.SelfClosingTagToken || tt == html.StartTagToken {
			switch (nam) {
			case "a":
				for {
					n, v, m := tok.TagAttr()
					if string(n) == "href" {
						TxtContent += string(v) + " "
						break
					}
					if !m {
						break
					}
				}
			case "li":
				TxtContent += "\n* "
			case "p":
				TxtContent += "\n"
			case "br":
				TxtContent += "\n"
			}
		}
	}
	return TxtContent
}

func main() {
	opt_rev := flag.Bool("r", false, "Reverse output")
	opt_html := flag.Bool("h", false, "Decode html")
	flag.Parse()
	if len(flag.Args()) < 1 {
		os.Exit(0)
	}
	resp, err := http.Get(flag.Args()[0])

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	rss := Rss{}
	decoder := xml.NewDecoder(resp.Body)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&rss)
	if err != nil {
		log.Fatal(err)
	}
	num := 0
	if *opt_rev {
		for k, _ := range rss.Channel.Items {
			rss.Channel.Items[k].num = k
		}
		sort.Slice(rss.Channel.Items, func(i, j int) bool {
			return rss.Channel.Items[i].num > rss.Channel.Items[j].num })
	}
	if rss.Channel.Title != "" {
		fmt.Printf("# %s\n\n", rss.Channel.Title);
	}
	for _, v := range rss.Channel.Items {
		if v.Link != "" {
			if *opt_html {
				v.Desc = html_decode(v.Desc)
			}
			num++
			fmt.Printf("## %s [%d]\n%s\n=> %s [%d]\n\n",
				v.Title, num, strings.TrimSpace(v.Desc)+"\n", v.Link, num)
		}
	}
}

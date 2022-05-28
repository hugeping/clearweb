package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

type Item struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Desc  string `xml:"description"`
	Full  string `xml:"encoded"`
	Date  string `xml:"pubDate"`
	num   int
	time  int64
}

type Channel struct {
	Items []Item `xml:"item"`
	Title string `xml:"title"`
	Desc  string `xml:"description"`
}

type Rss struct {
	Channel Channel `xml:"channel"`
}

func html_decode(s string) string {
	var links []string
	TxtContent := ""
	tok := html.NewTokenizer(strings.NewReader(s))
	prev := tok.Token()

	for {
		tt := tok.Next()
		if tt == html.ErrorToken {
			break
		}
		n, _ := tok.TagName()
		nam := string(n)
		switch {
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
			switch nam {
			case "a":
				for {
					n, v, m := tok.TagAttr()
					if string(n) == "href" && !strings.HasPrefix(string(v), "/") {
						links = append(links, string(v))
						TxtContent += fmt.Sprintf("[%c]", 64+len(links))
						break
					}
					if !m {
						break
					}
				}
			case "div":
				TxtContent += "\n\n"
			case "li":
				TxtContent += "\n* "
			case "p":
				TxtContent += "\n\n"
			case "br":
				TxtContent += "\n"
			}
		}
	}
	if len(links) > 0 {
		TxtContent += "\n"
	}
	for k, v := range links {
		TxtContent += fmt.Sprintf("=> %s [%c]\n", v, 65+k)
	}
	return TxtContent
}

func main() {
	opt_rev := flag.Bool("r", false, "Reverse output")
	opt_date := flag.Bool("d", false, "Sort by pubDate")
	opt_html := flag.Bool("h", false, "Decode html")
	opt_num := flag.Int("n", -1, "Limit")
	flag.Parse()
	rss := Rss{}
	var decoder *xml.Decoder
	if len(flag.Args()) < 1 {
		decoder = xml.NewDecoder(os.Stdin)
	} else {
		resp, err := http.Get(flag.Args()[0])
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		decoder = xml.NewDecoder(resp.Body)
	}
	decoder.CharsetReader = charset.NewReaderLabel
	for true {
		err := decoder.Decode(&rss)
		if err == nil {
			continue
		}
		if err == io.EOF {
			break
		}
		log.Fatal(err)
	}
	num := 0
	start := 0
	if *opt_rev || *opt_date {
		for k, item := range rss.Channel.Items {
			rss.Channel.Items[k].num = k
			if !*opt_date {
				continue
			}
			d := strings.TrimSpace(item.Date)
			t, err := time.Parse(time.RFC1123Z, d)
			if err != nil {
				t, err = time.Parse(time.RFC822, d)
			}
			if err == nil {
				rss.Channel.Items[k].time = t.Unix()
			}
		}
		sort.Slice(rss.Channel.Items, func(i, j int) bool {
			if *opt_date {
				if *opt_rev {
					return rss.Channel.Items[i].time < rss.Channel.Items[j].time
				}
				return rss.Channel.Items[i].time > rss.Channel.Items[j].time
			}
			return rss.Channel.Items[i].num > rss.Channel.Items[j].num
		})
		if *opt_num >= 0 {
			start = len(rss.Channel.Items) - *opt_num
		}
	}
	if rss.Channel.Title != "" {
		fmt.Printf("# %s\n\n", rss.Channel.Title)
	}
	for k, v := range rss.Channel.Items {
		if !*opt_rev {
			if *opt_num == 0 {
				break
			}
			*opt_num--
		} else {
			if k < start {
				continue
			}
		}
		if v.Link != "" {
			if *opt_html {
				v.Desc = html_decode(v.Desc)
			}
			num++
			dsc := strings.TrimSpace(v.Desc)
			if dsc != "" {
				dsc = "\n\n" + dsc + "\n"
			}
			fmt.Printf("## %s [%d]%s\n=> %s [%d]\n\n",
				v.Title, num, dsc, v.Link, num)
		}
	}
}

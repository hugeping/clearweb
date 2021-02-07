package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

type Item struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Desc  string `xml:"description"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type Rss struct {
	Channel Channel `xml:"channel"`
}

func main() {
	var page_opt *int = flag.Int("p", 1, "Page number")
	var server_opt *string = flag.String("s", "https://search.fedi.life", "Server url")

	flag.Parse()
	if len(flag.Args()) == 0 {
		os.Exit(0)
	}
	q := ""
	for _, v := range flag.Args() {
		if q != "" {
			q += " "
		}
		q += v
	}
	data := url.Values{
		"pageno":           {fmt.Sprintf("%d", *page_opt)},
		"category_general": {"1"},
		"format":           {"rss"},
		"q":                {q},
	}

	resp, err := http.PostForm(*server_opt+"/search", data)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	rss := Rss{}
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&rss)
	if err != nil {
		log.Fatal(err)
	}
	num := 0
	for _, v := range rss.Channel.Items {
		if v.Link != "" {
			if v.Desc != "" {
				v.Desc = v.Desc + "\n"
			}
			num++
			fmt.Printf("%s [%d]\n%s\n=> %s [%d]\n\n",
				v.Title, num, v.Desc, v.Link, num)
		}
	}
}

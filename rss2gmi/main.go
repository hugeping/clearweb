package main

import (
	"encoding/xml"
	"golang.org/x/net/html/charset"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
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

func main() {
	opt_rev := flag.Bool("r", false, "Reverse output")
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
			if v.Desc != "" {
				v.Desc = v.Desc + "\n"
			}
			num++
			fmt.Printf("%s [%d]\n%s\n=> %s [%d]\n\n",
				v.Title, num, v.Desc, v.Link, num)
		}
	}
}

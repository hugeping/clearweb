package main

import (
	"encoding/xml"
	"golang.org/x/net/html/charset"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Item struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Desc  string `xml:"description"`
	Full  string `xml:"encoded"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type Rss struct {
	Channel Channel `xml:"channel"`
}

func main() {
	if len(os.Args) < 2 {
		os.Exit(0)
	}
	resp, err := http.Get(os.Args[1])

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

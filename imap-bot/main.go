package main

import (
	"flag"
	"fmt"
	imap "github.com/emersion/go-imap"
	client "github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/k3a/html2text"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	limit := flag.Int("n", 1, "Number of last messages")
	server := flag.String("server", "127.0.0.1:993", "IMAP server")
	login := flag.String("login", "", "IMAP login")
	inbox := flag.String("mbox", "INBOX", "IMAP directory")
	password := flag.String("password", "", "IMAP password")
	id := flag.String("id", "", "Get body for message id")
	flag.Parse()

	if *login == "" {
		log.Fatal("No login")
	}

	if *password == "" {
		passwd, ok := os.LookupEnv("PASSWORD")
		if ok {
			*password = passwd
		}
	}

	c, err := client.DialTLS(*server, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Logout()

	if err := c.Login(*login, *password); err != nil {
		log.Fatal(err)
	}

	mailbox, err := c.Select(*inbox, false)
	if err != nil {
		log.Fatal(err)
	}
	seqSet := new(imap.SeqSet)
	seqSet.AddRange(mailbox.Messages-(uint32)(*limit)+1, mailbox.Messages)

	messages := make(chan *imap.Message, 10)
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem(), imap.FetchEnvelope}

	go func() {
		if err := c.Fetch(seqSet, items, messages); err != nil {
			log.Fatal(err)
		}
	}()

	for msg := range messages {
		if *id == "" {
			fmt.Printf("%s: %s\n",
				msg.Envelope.MessageId,
				msg.Envelope.Subject)
		} else if msg.Envelope.MessageId == *id {
			b := msg.GetBody(section)
			mr, err := mail.CreateReader(b)
			if err != nil {
				log.Fatal(err)
			}
			defer mr.Close()
			for {
				p, err := mr.NextPart()
				if err == io.EOF {
					break
				} else if err != nil {
					log.Fatal(err)
				}

				switch p.Header.(type) {
				case *mail.InlineHeader:
					t := p.Header.Get("Content-Type")
					b, err := ioutil.ReadAll(p.Body)
					if err != nil {
						continue
					}
					plain := string(b)
					if strings.Contains(t, "/html;") {
						plain = html2text.HTML2Text(string(b))
					}
					plain = strings.ReplaceAll(plain, "\r", "")
					if len(msg.Envelope.From) >= 1 {
						fmt.Println(msg.Envelope.From[0].PersonalName)
					}
					fmt.Println(msg.Envelope.Date)
					fmt.Println("")
					fmt.Println(plain)
					mr.Close()
					break
				}
			}
		}
	}
}

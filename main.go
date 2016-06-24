package main

import (
	"io/ioutil"
	"log"
	"strings"

	gzb "github.com/ifo/gozulipbot"
)

var emoji Set

func main() {
	emailAddress, apiKey, err := gzb.GetConfigFromFlags()
	if err != nil {
		log.Fatalln(err)
	}

	bot := gzb.Bot{
		Email:  emailAddress,
		APIKey: apiKey,
	}

	bot.Init()

	q, err := bot.RegisterAt()
	if err != nil {
		log.Fatal(err)
	}

	// load emoji
	emoji = getEmojiSet("emoji.txt")

	q.EventsCallback(respondToMessage)

	stop := make(chan struct{})
	<-stop
}

func respondToMessage(em gzb.EventMessage, err error) {
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("message received: " + em.Content)

	parts := strings.Fields(em.Content)

	if len(parts) < 2 {
		log.Println("invalid message")
		em.Queue.Bot.Respond(em, `ʕノ•ᴥ•ʔノ ︵ ┻━┻`)
		return
	}

	emj := parts[len(parts)-1]

	if emoji.Has(emj) {
		em.Queue.Bot.Respond(em, ":"+emj+":\n:octopus:")
	} else {
		log.Println("invalid emoji " + emj)
		em.Queue.Bot.Respond(em, `¯\_(ツ)_/¯`)
	}
}

type Set map[string]struct{}

func getEmojiSet(fileName string) Set {
	ebts, _ := ioutil.ReadFile(fileName)
	out := Set{}
	for _, e := range strings.Fields(string(ebts)) {
		out[e] = struct{}{}
	}

	return out
}

func (s *Set) Has(elem string) bool {
	_, ok := (*s)[elem]
	return ok
}

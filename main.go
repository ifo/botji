package main

import (
	"io/ioutil"
	"log"
	"strings"

	gzb "github.com/ifo/gozulipbot"
)

var emoji Set
var bases Set

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
	bases = getEmojiSet("bases.txt")

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

	fields := strings.Fields(em.Content)

	if len(fields) < 2 {
		log.Println("invalid message")
		em.Queue.Bot.Respond(em, `ʕノ•ᴥ•ʔノ ︵ ┻━┻`)
		return
	}

	shrug := true
	base := "octopus"
	top := ""

	for _, w := range fields {
		if bases.Has(w) {
			base = w
			shrug = false
		} else if emoji.Has(w) {
			top = w
			shrug = false
		}
	}

	// no emoji found, shrug
	if shrug {
		log.Println("invalid emoji " + top)
		em.Queue.Bot.Respond(em, `¯\_(ツ)_/¯`)
		return
	}

	// no top, use whatever base was found as top
	if top == "" {
		top = base
		base = "octopus"
	}

	em.Queue.Bot.Respond(em, ":"+top+":\n:"+base+":")
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

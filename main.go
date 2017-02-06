package main

import (
	"log"
	"os"
	"strings"

	gzb "github.com/ifo/gozulipbot"
)

var emoji Set
var bases Set

func main() {
	bot := gzb.Bot{}
	err := bot.GetConfigFromFlags()
	if err != nil {
		log.Fatalln(err)
	}
	bot.Init()

	q, err := bot.RegisterAt()
	if err != nil {
		log.Fatal(err)
	}

	// load emoji
	emoji = getEmojiSet("emoji.txt")
	bases = getEmojiSet("bases.txt")

	// setup log file
	f, err := os.OpenFile("botji.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	log.SetOutput(f)

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

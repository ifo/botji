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

	// Load emoji.
	emoji = getEmojiSet("emoji.txt")
	realm, err := bot.RealmEmojiSet()
	if err != nil {
		log.Fatal(err)
	}
	emoji.Union(realm)

	// Setup log file.
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

	unbelievable := true
	emjs := []string{}
	for _, w := range fields {
		if emoji.Has(w) {
			emjs = append(emjs, w)
			unbelievable = false
		}
	}

	// No emoji found, this is unbelievable.
	if unbelievable {
		em.Queue.Bot.React(em, "astonished")
		return
	}

	// Send reactions!
	for _, e := range emjs {
		em.Queue.Bot.React(em, e)
	}
}

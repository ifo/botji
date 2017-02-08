package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	gzb "github.com/ifo/gozulipbot"
)

var emoji Set

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

	// Stay subscribed to every stream.
	go subscribeEveryDay(bot)

	q.EventsCallback(reactToMessage)

	stop := make(chan struct{})
	<-stop
}

func reactToMessage(em gzb.EventMessage, err error) {
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

	// No emoji found; this is unbelievable.
	if unbelievable {
		em.Queue.Bot.React(em, "astonished")
		return
	}

	// Send reactions!
	for _, e := range emjs {
		em.Queue.Bot.React(em, e)
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

func (s *Set) Union(s2 Set) {
	for k := range s2 {
		(*s)[k] = struct{}{}
	}
}

func subscribeEveryDay(bot gzb.Bot) {
	for {
		streams, err := bot.GetStreams()
		if err != nil {
			log.Println(err)
		}
		resp, err := bot.Subscribe(streams)
		if err != nil {
			log.Println(err)
		}
		if resp.StatusCode >= 400 {
			log.Println(fmt.Errorf("Subscribe got error code %d", resp.StatusCode))
		}
		time.Sleep(24 * time.Hour)
	}
}

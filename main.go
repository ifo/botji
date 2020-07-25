package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode"

	gzb "github.com/ifo/gozulipbot"
)

var emoji Set

func main() {
	logFile := flag.String("logfile", "", "log file location, defaults to standard out")

	bot := gzb.Bot{}
	// This calls flag.Parse()
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

	if *logFile != "" {
		// Setup log file.
		f, err := os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

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

	emjs := parseEmoji(em.Content)
	unbelievable := false
	if len(emjs) == 0 {
		unbelievable = true
	}

	// No emoji found; this is unbelievable.
	if unbelievable {
		em.Queue.Bot.React(em, "astonished")
		return
	}

	// Send reactions!
	for e := range emjs {
		em.Queue.Bot.React(em, e)
	}
}

func parseEmoji(msg string) Set {
	out := Set{}
	clean := func(r rune) rune {
		switch {
		case r == ':' || r == '_' || r == '-':
			return ' '
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			return -1
		default:
			return r
		}
	}
	fields := strings.Fields(strings.ToLower(strings.Map(clean, msg)))

	for i := range fields {
		length := 8
		if len(fields)-i < length {
			length = len(fields) - i
		}
		for j := length + i; i < j; j-- {
			emj := strings.Join(fields[i:j], "_")
			if emoji.Has(emj) {
				out[emj] = struct{}{}
				// We've found an emoji, move the cursor to the end of this group.
				i = j
				break
			}
		}
	}

	return out
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

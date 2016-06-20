package main

import (
	"log"
	"strings"

	gzb "github.com/ifo/gozulipbot"
)

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
		em.Queue.Bot.Respond(em, `¯\_(ツ)_/¯`)
		return
	}

	emoji := parts[len(parts)-1]

	em.Queue.Bot.Respond(em, ":"+emoji+":\n:octopus:")
}

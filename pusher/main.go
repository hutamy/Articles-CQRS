package main

import (
	"articles/event"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	NatsAddress string `envconfig:"NATS_ADDRESS"`
}

func main() {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}

	hub := HubInit()
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		es, err := event.NatsInit(fmt.Sprintf("nats://%s", config.NatsAddress))
		if err != nil {
			log.Println(err)
			return err
		}

		err = es.OnArticleCreated(func(m event.ArticleCreatedMessage) {
			log.Printf("Article received: %v\n", m)
			hub.broadcast(newArticleCreatedMessage(m.ID, m.Author, m.Title, m.Body, m.Created), nil)
		})
		if err != nil {
			log.Println(err)
			return err
		}

		event.SetEventStore(es)
		return nil
	})
	defer event.Close()

	go hub.run()
	http.HandleFunc("/pusher", hub.handleWebSocket)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

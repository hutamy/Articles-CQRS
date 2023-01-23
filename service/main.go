package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"

	"articles/db"
	"articles/event"
	"articles/util"
)

type Config struct {
	PostgresDB       string `envconfig:"POSTGRES_DB"`
	PostgresUser     string `envconfig:"POSTGRES_USER"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress      string `envconfig:"NATS_ADDRESS"`
}

func newRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.
		HandleFunc("/articles", createArticleHandler).
		Methods(http.MethodPost)
	router.Use(mux.CORSMethodMiddleware(router))
	return
}

func main() {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}

	util.Retry(func(_ int) error {
		connectionStr := fmt.Sprintf(
			"postgres://%s:%s@postgres/%s?sslmode=disable",
			config.PostgresUser,
			config.PostgresPassword,
			config.PostgresDB,
		)
		repo, err := db.PostgresInit(connectionStr)
		if err != nil {
			log.Println(err)
			return err
		}
		db.SetRepository(repo)
		return nil
	})
	defer db.Close()

	util.Retry(func(_ int) error {
		connectionStr := fmt.Sprintf("nats://%s", config.NatsAddress)
		es, err := event.NatsInit(connectionStr)
		if err != nil {
			log.Println(err)
			return err
		}
		event.SetEventStore(es)
		return nil
	})
	defer event.Close()

	router := newRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

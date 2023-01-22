package main

import (
	"articles/db"
	"articles/event"
	"articles/search"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	PostgresDB           string `envconfig:"POSTGRES_DB"`
	PostgresUser         string `envconfig:"POSTGRES_USER"`
	PostgresPassword     string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress          string `envconfig:"NATS_ADDRESS"`
	ElasticsearchAddress string `envconfig:"ELASTICSEARCH_ADDRESS"`
}

func newRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/articles", listArticlesHandler).Methods("GET")
	return
}

func main() {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}

	retry.ForeverSleep(2*time.Second, func(_ int) error {
		addr := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", config.PostgresUser, config.PostgresPassword, config.PostgresDB)
		repo, err := db.PostgresInit(addr)
		if err != nil {
			log.Println(err)
			return err
		}
		db.SetRepository(repo)
		return nil
	})
	defer db.Close()

	retry.ForeverSleep(2*time.Second, func(_ int) error {
		es, err := search.ElasticInit(fmt.Sprintf("http://%s", config.ElasticsearchAddress))
		if err != nil {
			log.Println(err)
			return err
		}
		search.SetRepository(es)
		return nil
	})
	defer search.Close()

	retry.ForeverSleep(2*time.Second, func(_ int) error {
		es, err := event.NatsInit(fmt.Sprintf("nats://%s", config.NatsAddress))
		if err != nil {
			log.Println(err)
			return err
		}
		err = es.OnArticleCreated(onArticleCreated)
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

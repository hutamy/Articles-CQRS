package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"

	"articles/db"
	"articles/event"
	"articles/search"
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
	router.
		HandleFunc("/articles", listArticlesHandler).
		Methods(http.MethodGet)
	router.Use(mux.CORSMethodMiddleware(router))
	return
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	retry.ForeverSleep(2*time.Second, func(_ int) error {
		addr := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)
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
		addr := fmt.Sprintf("http://%s", cfg.ElasticsearchAddress)
		es, err := search.ElasticInit(addr)
		if err != nil {
			log.Println(err)
			return err
		}
		search.SetRepository(es)
		return nil
	})
	defer search.Close()

	retry.ForeverSleep(2*time.Second, func(_ int) error {
		addr := fmt.Sprintf("nats://%s", cfg.NatsAddress)
		es, err := event.NatsInit(addr)
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

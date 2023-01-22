package main

import (
	"articles/db"
	"articles/event"
	"articles/schema"
	"articles/util"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type request struct {
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func createArticleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body request
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	article := schema.Article{
		ID:      0,
		Author:  body.Author,
		Title:   body.Title,
		Body:    body.Body,
		Created: time.Now().UTC(),
	}

	respArticle, err := db.InsertArticle(ctx, article)
	if err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Failed to create article")
	}

	if err := event.PublishArticleCreated(respArticle); err != nil {
		log.Println(err)
	}

	util.ResponseOk(w, respArticle)
}

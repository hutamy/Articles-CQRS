package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"articles/db"
	"articles/event"
	"articles/schema"
	"articles/util"

	"github.com/segmentio/ksuid"
)

type response struct {
	ID string `json:"id"`
}

type request struct {
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func createArticleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body request
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.ResponseError(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	err = json.Unmarshal(b, &body)
	if err != nil {
		util.ResponseError(w, http.StatusBadRequest, "Failed to unmarshall request body")
		return
	}

	if err := checkEmptyBody(body); err != "" {
		util.ResponseError(w, http.StatusBadRequest, err)
		return
	}

	created := time.Now().UTC()
	id, err := ksuid.NewRandomWithTime(created)
	if err != nil {
		util.ResponseError(w, http.StatusInternalServerError, "Failed to create article")
		return
	}
	article := schema.Article{
		ID:      id.String(),
		Author:  body.Author,
		Title:   body.Title,
		Body:    body.Body,
		Created: created,
	}
	if err := db.InsertArticle(ctx, article); err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Failed to create article")
		return
	}
	if err := event.PublishArticleCreated(article); err != nil {
		log.Println(err)
	}
	util.ResponseOk(w, response{ID: article.ID})
}

func checkEmptyBody(body request) string {
	if body.Author == "" {
		return "Author is required"
	}
	if body.Title == "" {
		return "Title is required"
	}
	if body.Body == "" {
		return "Body is required"
	}
	return ""
}

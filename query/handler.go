package main

import (
	"articles/db"
	"articles/event"
	"articles/schema"
	"articles/search"
	"articles/util"
	"context"
	"log"
	"net/http"
)

func onArticleCreated(m event.ArticleCreatedMessage) {
	article := schema.Article{
		ID:      m.ID,
		Author:  m.Author,
		Title:   m.Title,
		Body:    m.Body,
		Created: m.Created,
	}
	if err := search.InsertArticle(context.Background(), article); err != nil {
		log.Println(err)
	}
}

func listArticlesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	articles, err := db.ListArticles(ctx)
	if err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Could not fetch articles")
		return
	}
	util.ResponseOk(w, articles)
}

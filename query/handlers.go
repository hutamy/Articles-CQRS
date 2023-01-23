package main

import (
	"context"
	"log"
	"net/http"

	"articles/event"
	"articles/schema"
	"articles/search"
	"articles/util"
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
	var err error
	ctx := r.Context()
	query := r.FormValue("query")
	author := r.FormValue("author")
	articles, err := search.ListArticles(ctx, query, author)
	if err != nil {
		log.Println(err)
		util.ResponseOk(w, []schema.Article{})
		return
	}
	util.ResponseOk(w, articles)
}

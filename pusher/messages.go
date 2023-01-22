package main

import "time"

const (
	KindArticleCreated = iota + 1
)

type ArticleCreatedMessage struct {
	Kind    uint32    `json:"kind"`
	ID      int       `json:"id"`
	Author  string    `json:"author"`
	Title   string    `json:"title"`
	Body    string    `json:"body"`
	Created time.Time `json:"created"`
}

func newArticleCreatedMessage(id int, author string, title string, body string, created time.Time) *ArticleCreatedMessage {
	return &ArticleCreatedMessage{
		Kind:    KindArticleCreated,
		ID:      id,
		Author:  author,
		Title:   title,
		Body:    body,
		Created: created,
	}
}

package event

import "time"

type Message interface {
	Key() string
}

type ArticleCreatedMessage struct {
	ID      int
	Author  string
	Title   string
	Body    string
	Created time.Time
}

func (a *ArticleCreatedMessage) Key() string {
	return "article.created"
}

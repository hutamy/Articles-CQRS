package event

import (
	"time"
)

type Message interface {
	Key() string
}

type ArticleCreatedMessage struct {
	ID      string
	Author  string
	Title   string
	Body    string
	Created time.Time
}

func (m *ArticleCreatedMessage) Key() string {
	return "article.created"
}

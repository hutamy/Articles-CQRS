package event

import "articles/schema"

type Store interface {
	Close()
	PublishArticleCreated(article schema.Article) error
	SubscribeArticleCreated() (<-chan ArticleCreatedMessage, error)
	OnArticleCreated(f func(ArticleCreatedMessage)) error
}

var s Store

func SetEventStore(es Store) {
	s = es
}

func Close() {
	s.Close()
}

func PublishArticleCreated(article schema.Article) error {
	return s.PublishArticleCreated(article)
}

func SubscribeArticleCreated() (<-chan ArticleCreatedMessage, error) {
	return s.SubscribeArticleCreated()
}

func OnArticleCreated(f func(ArticleCreatedMessage)) error {
	return s.OnArticleCreated(f)
}

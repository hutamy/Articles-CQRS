package event

import "articles/schema"

type EventStore interface {
	Close()
	PublishArticleCreated(article schema.Article) error
	SubscribeArticleCreated() (<-chan ArticleCreatedMessage, error)
	OnArticleCreated(f func(ArticleCreatedMessage)) error
}

var impl EventStore

func SetEventStore(es EventStore) {
	impl = es
}

func Close() {
	impl.Close()
}

func PublishArticleCreated(article schema.Article) error {
	return impl.PublishArticleCreated(article)
}

func SubscribeArticleCreated() (<-chan ArticleCreatedMessage, error) {
	return impl.SubscribeArticleCreated()
}

func OnArticleCreated(f func(ArticleCreatedMessage)) error {
	return impl.OnArticleCreated(f)
}

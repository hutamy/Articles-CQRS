package search

import (
	"articles/schema"
	"context"
)

type Repository interface {
	Close()
	InsertArticle(ctx context.Context, article schema.Article) error
	SearchArticles(ctx context.Context, query string) ([]schema.Article, error)
}

var impl Repository

func SetRepository(repository Repository) {
	impl = repository
}

func Close() {
	impl.Close()
}

func InsertArticle(ctx context.Context, article schema.Article) error {
	return impl.InsertArticle(ctx, article)
}

func SearchArticles(ctx context.Context, query string) ([]schema.Article, error) {
	return impl.SearchArticles(ctx, query)
}

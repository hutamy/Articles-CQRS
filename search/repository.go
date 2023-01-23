package search

import (
	"context"

	"articles/schema"
)

type Repository interface {
	Close()
	InsertArticle(ctx context.Context, article schema.Article) error
	ListArticles(ctx context.Context, query string, author string) ([]schema.Article, error)
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

func ListArticles(ctx context.Context, query string, author string) ([]schema.Article, error) {
	return impl.ListArticles(ctx, query, author)
}

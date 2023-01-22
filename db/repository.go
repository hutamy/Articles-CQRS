package db

import (
	"articles/schema"
	"context"
)

type Repository interface {
	Close()
	InsertArticle(ctx context.Context, article schema.Article) (schema.Article, error)
	ListArticles(ctx context.Context) ([]schema.Article, error)
}

var impl Repository

func SetRepository(repository Repository) {
	impl = repository
}

func Close() {
	impl.Close()
}

func InsertArticle(ctx context.Context, article schema.Article) (schema.Article, error) {
	return impl.InsertArticle(ctx, article)
}

func ListArticles(ctx context.Context) ([]schema.Article, error) {
	return impl.ListArticles(ctx)
}

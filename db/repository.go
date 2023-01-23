package db

import (
	"context"

	"articles/schema"
)

type Repository interface {
	Close()
	InsertArticle(ctx context.Context, article schema.Article) error
}

var r Repository

func SetRepository(repository Repository) {
	r = repository
}

func Close() {
	r.Close()
}

func InsertArticle(ctx context.Context, article schema.Article) error {
	return r.InsertArticle(ctx, article)
}

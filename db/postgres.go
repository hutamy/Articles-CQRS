package db

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"articles/schema"
)

type PostgresRepository struct {
	db *sql.DB
}

func PostgresInit(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{
		db,
	}, nil
}

func (r *PostgresRepository) Close() {
	if err := r.db.Close(); err != nil {
		log.Fatal(err)
	}
}

func (r *PostgresRepository) InsertArticle(ctx context.Context, article schema.Article) error {
	_, err := r.db.Exec(
		"INSERT INTO articles(id, author, title, body, created) VALUES($1, $2, $3, $4, $5)",
		article.ID,
		article.Author,
		article.Title,
		article.Body,
		article.Created,
	)
	return err
}

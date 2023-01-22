package db

import (
	"articles/schema"
	"context"
	"database/sql"
	"log"
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

func (r *PostgresRepository) InsertArticle(ctx context.Context, article schema.Article) (schema.Article, error) {
	row := r.db.QueryRow(
		`insert into articles (author, title, body, created) values ($1, $2, $3, $4) returning *;`,
		article.Author,
		article.Title,
		article.Body,
		article.Created,
	)
	var result schema.Article
	if err := row.Scan(&result); err != nil {
		log.Fatal(err)
		return schema.Article{}, err
	}
	return result, nil
}

func (r *PostgresRepository) ListArticles(ctx context.Context) ([]schema.Article, error) {
	rows, err := r.db.Query(`select * from articles order by id desc;`)
	if err != nil {
		return []schema.Article{}, err
	}
	defer rows.Close()

	articles := []schema.Article{}
	for rows.Next() {
		row := schema.Article{}
		if err = rows.Scan(&row); err != nil {
			articles = append(articles, row)
		}
	}

	if err = rows.Err(); err != nil {
		return []schema.Article{}, err
	}

	return articles, nil
}

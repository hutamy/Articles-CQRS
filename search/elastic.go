package search

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	elastic "github.com/elastic/go-elasticsearch/v7"

	"articles/schema"
)

type ElasticRepository struct {
	client *elastic.Client
}

func ElasticInit(url string) (*ElasticRepository, error) {
	client, err := elastic.NewClient(elastic.Config{
		Addresses: []string{url},
	})
	if err != nil {
		return nil, err
	}
	_, err = client.Info()
	if err != nil {
		return nil, err
	}
	return &ElasticRepository{client}, nil
}

func (r *ElasticRepository) Close() {
}

func (r *ElasticRepository) InsertArticle(ctx context.Context, article schema.Article) error {
	body, _ := json.Marshal(article)
	_, err := r.client.Index(
		"articles",
		bytes.NewReader(body),
		r.client.Index.WithDocumentID(article.ID),
		r.client.Index.WithRefresh("wait_for"),
	)
	return err
}

func (r *ElasticRepository) ListArticles(ctx context.Context, query string, author string) (result []schema.Article, err error) {
	var buf bytes.Buffer
	reqBody := map[string]interface{}{}
	multiMatch := []map[string]interface{}{}

	if query != "" {
		multiMatchBodyTitle := map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":            query,
				"fields":           []string{"body", "title"},
				"fuzziness":        3,
				"cutoff_frequency": 0.0001,
			},
		}
		multiMatch = append(multiMatch, multiMatchBodyTitle)
	}

	if author != "" {
		matchAuthor := map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":            author,
				"fields":           []string{"author"},
				"fuzziness":        3,
				"cutoff_frequency": 0.0001,
			},
		}
		multiMatch = append(multiMatch, matchAuthor)
	}

	if query != "" || author != "" {
		reqBody["query"] = map[string]interface{}{
			"dis_max": map[string]interface{}{
				"queries": multiMatch,
			},
		}
	}

	if err = json.NewEncoder(&buf).Encode(reqBody); err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex("articles"),
		r.client.Search.WithBody(&buf),
		r.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			result = nil
		}
	}()
	if res.IsError() {
		return nil, errors.New("search failed")
	}

	type Response struct {
		Took int64
		Hits struct {
			Total struct {
				Value int64
			}
			Hits []*struct {
				Source schema.Article `json:"_source"`
			}
		}
	}
	resBody := Response{}
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return nil, err
	}

	var articles []schema.Article
	for _, hit := range resBody.Hits.Hits {
		articles = append(articles, hit.Source)
	}
	return articles, nil
}

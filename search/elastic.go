package search

import (
	"articles/schema"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strconv"

	elastic "github.com/elastic/go-elasticsearch/v7"
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
	return &ElasticRepository{client: client}, nil
}

func (r *ElasticRepository) Close() {}

func (r *ElasticRepository) InsertArticle(ctx context.Context, article schema.Article) error {
	strId := strconv.Itoa(article.ID)
	body, _ := json.Marshal(article)
	_, err := r.client.Index(
		"articles",
		bytes.NewReader(body),
		r.client.Index.WithDocumentID(strId),
		r.client.Index.WithRefresh("wait_for"),
	)
	return err
}

func (r *ElasticRepository) SearchArticles(ctx context.Context, query string) (result []schema.Article, err error) {
	var buf bytes.Buffer
	reqBody := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":            query,
				"fields":           []string{"body"},
				"fuzziness":        3,
				"cutoff_frequency": 0.0001,
			},
		},
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

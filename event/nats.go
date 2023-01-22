package event

import (
	"articles/schema"
	"bytes"
	"encoding/gob"
	"log"

	"github.com/nats-io/nats.go"
)

type NatsEventStore struct {
	nc                         *nats.Conn
	articleCreatedSubscription *nats.Subscription
	articleCreatedChan         chan ArticleCreatedMessage
}

func NatsInit(url string) (*NatsEventStore, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsEventStore{nc: nc}, nil
}

func (es *NatsEventStore) SubscribeArticleCreated() (<-chan ArticleCreatedMessage, error) {
	m := ArticleCreatedMessage{}
	es.articleCreatedChan = make(chan ArticleCreatedMessage, 64)
	ch := make(chan *nats.Msg, 64)
	var err error
	es.articleCreatedSubscription, err = es.nc.ChanSubscribe(m.Key(), ch)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case msg := <-ch:
				if err := es.readMessage(msg.Data, &m); err != nil {
					log.Fatal(err)
				}
				es.articleCreatedChan <- m
			}
		}
	}()
	return (<-chan ArticleCreatedMessage)(es.articleCreatedChan), nil
}

func (es *NatsEventStore) OnArticleCreated(f func(ArticleCreatedMessage)) (err error) {
	m := ArticleCreatedMessage{}
	es.articleCreatedSubscription, err = es.nc.Subscribe(m.Key(), func(msg *nats.Msg) {
		es.readMessage(msg.Data, &m)
		f(m)
	})
	return
}

func (es *NatsEventStore) Close() {
	if es.nc != nil {
		es.nc.Close()
	}

	if es.articleCreatedSubscription != nil {
		es.articleCreatedSubscription.Unsubscribe()
	}

	close(es.articleCreatedChan)
}

func (es *NatsEventStore) PublishArticleCreated(article schema.Article) error {
	m := ArticleCreatedMessage{
		article.ID,
		article.Author,
		article.Title,
		article.Body,
		article.Created,
	}
	data, err := es.writeMessage(&m)
	if err != nil {
		return err
	}
	return es.nc.Publish(m.Key(), data)
}

func (es *NatsEventStore) writeMessage(m Message) ([]byte, error) {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(m)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (es *NatsEventStore) readMessage(data []byte, m interface{}) error {
	b := bytes.Buffer{}
	b.Write(data)
	return gob.NewDecoder(&b).Decode(m)
}

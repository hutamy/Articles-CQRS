# Articles

## Prerequisites
- Golang 
- Docker 
- Makefile

## Tech Stack
- CQRS Pattern
- PostgreSQL
- NATS
- Elasticsearch
- Docker

## Run Application
- make up or docker-compose up -d --build

## Directory structure
```
├── db
│   └── postgres.go
│   └── repository.go
├── event
│   └── event.go
│   └── messages.go
│   └── nats.go
├── nginx
│   └── Dockerfile
│   └── nginx.conf
├── postgres
│   └── Dockerfile
│   └── up.sql
├── query
│   └── handler.go
│   └── main.go
├── schema
│   └── model.go
├── search
│   └── elastic.go
│   └── repository.go
├── service
│   └── handler.go
│   └── main.go
├── util
│   └── util.go
├── vendor
├── docker-compose.yaml
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
├── README.md
```

## API Reference

#### Create an article

```http
  POST /articles
```
Request body:
| Field           | Type     | Description                     |
| :-------------- | :------- | :------------------------------ |
| `author`        | `string` | **Required**. author of article |
| `title`         | `string` | **Required**. title of article  |
| `body`          | `string` | **Required**. body of article   |

Examples:
```bash
curl --location --request POST 'http://localhost:8080/articles' \
--header 'Content-Type: application/json' \
--data-raw '{
    "author": "hutamy",
    "title": "Create CQRS Microservice",
    "body": "this is how to create microservice with cqrs pattern"
}'
```

Result
```
{
  "id": "2KjDVL85Og9R3FosuyotaaVxFqK"
}
```

#### Get articles

```http
  GET /articles
```
Request query:
| Field           | Type     | Description                                          |
| :-------------- | :------- | :--------------------------------------------------- |
| `author`        | `string` | **Not Required**. list of articles by author         |
| `query`         | `string` | **Not Required**. list of articles by title or body  |

Examples:
```bash
curl --location --request GET 'http://localhost:8080/articles?query=cqrs&author=hutamy'
```

Result
```
[
  {
    "id": "2KjDVL85Og9R3FosuyotaaVxFqK",
    "author": "hutamy",
    "title": "Create CQRS Microservice",
    "body": "this is how to create microservice with cqrs pattern",
    "created": "2023-01-23 01:10:44.29563+07"
  }
]
```

## Authors

- [@hutamy](https://www.github.com/hutamy)
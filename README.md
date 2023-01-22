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
- make up

## Directory structure
```
├── db
│   └── postgres.go
│   └── repository.go
├── event
│   └── event.go
│   └── messages.go
│   └── nats.go
├── postgres
│   └── Dockerfile
│   └── up.sql
├── pusher
│   └── client.go
│   └── hub.go
│   └── main.go
│   └── messgaes.go
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
├── docker-compose.yml
├── Dockerfile
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
curl -X POST "localhost:8080/articles" \
-H "Accept: application/json" \
-H "Content-Type: application/json" \
-d '{"author":"hutamy","title":"Create CQRS Microservice","body":"this is how to create microservice with cqrs pattern"}'
```

Result
```
{
  "id": 1
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
curl -X GET "localhost:8080/articles" \
-H "Accept: application/json" \
-H "Content-Type: application/json" \
```

Result
```
[
  {
    "id": 1,
    "author": "hutamy",
    "title": "Create CQRS Microservice",
    "body": "this is how to create microservice with cqrs pattern",
    "created": "2023-01-23 01:10:44.29563+07"
  }
]
```

## Authors

- [@hutamy](https://www.github.com/hutamy)
package schema

import (
	"time"
)

type Article struct {
	ID      string    `json:"id"`
	Author  string    `json:"author"`
	Title   string    `json:"title"`
	Body    string    `json:"body"`
	Created time.Time `json:"created"`
}

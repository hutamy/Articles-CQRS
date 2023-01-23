package util

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type RetryFunc func(int) error

func ResponseOk(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	log.Println(body)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Fatal(err)
	}
}

func ResponseError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	body := map[string]string{
		"error": message,
	}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Fatal(err)
	}
}

func Retry(f RetryFunc) {
	for i := 0; ; i++ {
		err := f(i)
		if err == nil {
			return
		}
		time.Sleep(2 * time.Second)
	}
}

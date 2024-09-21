package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var startedAt = time.Now()

func main() {

	http.HandleFunc("/healthz", Healthz)
	http.HandleFunc("/", Hello)
	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
	w.Write([]byte("Hello"))
	w.WriteHeader(http.StatusOK)
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	duration := time.Since(startedAt)

	if duration.Seconds() < 10 {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Duration: %v", duration.Seconds())))
	} else {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}
}

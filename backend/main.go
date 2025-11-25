package main

import (
	"log"
	"net/http"
)

func main() {
	client := NewOpenF1Client()
	router := NewRouter(client)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

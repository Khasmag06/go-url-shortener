package main

import (
	"github.com/Khasmag06/go-url-shortener/internal/app/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.Route)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

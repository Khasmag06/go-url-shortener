package main

import (
	"github.com/Khasmag06/go-url-shortener/internal/app/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.PostHandler)
	r.HandleFunc("/{id}", handlers.GetHandler).Methods(http.MethodGet)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

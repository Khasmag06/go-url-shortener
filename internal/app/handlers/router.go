package handlers

import (
	"net/http"
	"regexp"
)

var validPathPost = regexp.MustCompile(`^/$`)
var validPathGet = regexp.MustCompile(`^/([a-zA-Z]+)$`)

func Route(w http.ResponseWriter, r *http.Request) {
	switch {
	case validPathPost.MatchString(r.URL.Path):
		PostHandler(w, r)
	case validPathGet.MatchString(r.URL.Path):
		GetHandler(w, r)

	}
}

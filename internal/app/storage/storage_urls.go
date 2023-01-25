package storage

import (
	"errors"
)

type Storage interface {
	Add(shortURL *ShortURL)
	Get(id string) (*ShortURL, error)
}

type ShortURL struct {
	ID          string
	OriginalURL string
}

type URLStorage struct {
	urls []*ShortURL
}

func (u *URLStorage) Add(s *ShortURL) {
	u.urls = append(u.urls, s)

}

func (u *URLStorage) Get(id string) (*ShortURL, error) {
	for _, el := range u.urls {
		if el.ID == id {
			return el, nil
		}
	}
	return &ShortURL{"", ""}, errors.New("not found")
}

var short = ShortURL{"/google", "https://www.google.com/"}
var Urls Storage = &URLStorage{[]*ShortURL{&short}}

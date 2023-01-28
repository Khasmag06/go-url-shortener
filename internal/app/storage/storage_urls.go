package storage

import (
	"encoding/json"
	"errors"
	"github.com/Khasmag06/go-url-shortener/config"
	"os"
)

type Storage interface {
	Add(shortURL *ShortURL)
	Get(id string) (*ShortURL, error)
}

type ShortURL struct {
	ID          string `json:"id"`
	OriginalURL string `json:"originalURL"`
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

type URLStorageFile struct {
	file *os.File
}

func (u *URLStorageFile) Add(s *ShortURL) {
	file, _ := os.OpenFile(config.Cfg.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	u.file = file
	json.NewEncoder(u.file).Encode(&s)
	defer u.file.Close()
}

func (u *URLStorageFile) Get(id string) (*ShortURL, error) {
	file, _ := os.OpenFile(config.Cfg.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0777)
	u.file = file
	defer u.file.Close()
	short := &ShortURL{}
	encoder := json.NewDecoder(u.file)
	for {
		err := encoder.Decode(&short)
		if err != nil {
			return nil, err
		}
		if short.ID == id {
			return short, nil
		}
		if short.ID == "" {
			break
		}
	}
	return short, errors.New("not found")
}

func NewStorage() Storage {
	if config.Cfg.FileStoragePath != "" {
		return &URLStorageFile{}
	}
	var short = ShortURL{"/google", "https://www.google.com/"}

	return &URLStorage{[]*ShortURL{&short}}
}

var Urls = NewStorage()

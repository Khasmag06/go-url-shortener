package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var ErrNotFound = errors.New("not found")

type Storage interface {
	AddShortURL(shortURL *ShortURL) error
	GetShortURL(id string) (*ShortURL, error)
}

type ShortURL struct {
	ID          string `json:"id"`
	OriginalURL string `json:"originalURL"`
}

type URLStorage struct {
	urls []*ShortURL
}

func (u *URLStorage) AddShortURL(s *ShortURL) error {
	u.urls = append(u.urls, s)
	return nil

}

func (u *URLStorage) GetShortURL(id string) (*ShortURL, error) {
	for _, el := range u.urls {
		if el.ID == id {
			return el, nil
		}
	}
	return &ShortURL{"", ""}, ErrNotFound
}

type URLStorageFile struct {
	filepath string
}

func (u *URLStorageFile) AddShortURL(s *ShortURL) error {
	file, err := os.OpenFile(u.filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return fmt.Errorf("unable to open file: %w", err)
	}
	defer file.Close()
	json.NewEncoder(file).Encode(&s)
	return nil
}

func (u *URLStorageFile) GetShortURL(id string) (*ShortURL, error) {
	file, _ := os.OpenFile(u.filepath, os.O_RDONLY|os.O_CREATE, 0777)
	defer file.Close()
	short := &ShortURL{}
	encoder := json.NewDecoder(file)
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
	return short, ErrNotFound
}

func NewMemoryStorage() Storage {
	var short = &ShortURL{"/google", "https://www.google.com/"}
	return &URLStorage{
		urls: []*ShortURL{short},
	}
}

func NewFileStorage(filepath string) Storage {
	return &URLStorageFile{filepath: filepath}
}

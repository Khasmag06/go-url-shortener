package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

var ErrNotFound = errors.New("not found")

type Storage interface {
	AddShortURL(userID string, shortURL *ShortURL) error
	GetShortURL(short string) (*ShortURL, error)
	GetAllShortURL(userID string) []*ShortURL
}

type ShortURL struct {
	ID          string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type MemStorage struct {
	userURLs map[string][]*ShortURL
	urls     []*ShortURL
}

func (ms *MemStorage) AddShortURL(userID string, s *ShortURL) error {
	ms.userURLs[userID] = append(ms.userURLs[userID], s)
	ms.urls = append(ms.urls, s)
	return nil

}

func (ms *MemStorage) GetShortURL(id string) (*ShortURL, error) {
	for _, el := range ms.urls {
		if el.ID == id {
			return el, nil
		}
	}
	return nil, ErrNotFound
}

func (ms *MemStorage) GetAllShortURL(userID string) []*ShortURL {
	return ms.userURLs[userID]
}

func NewMemoryStorage() Storage {
	var short = &ShortURL{"google", "https://www.google.com/"}
	var userID = "12345"
	var shortsList = []*ShortURL{short}
	return &MemStorage{
		userURLs: map[string][]*ShortURL{userID: shortsList},
		urls:     shortsList,
	}
}

type FileStorage struct {
	*MemStorage
	f *os.File
}

func (fs *FileStorage) AddShortURL(userID string, s *ShortURL) error {
	if err := fs.MemStorage.AddShortURL(userID, s); err != nil {
		return fmt.Errorf("unable to add new key in memorystorage: %w", err)
	}
	err := fs.f.Truncate(0)
	if err != nil {
		return fmt.Errorf("unable to truncate file: %w", err)
	}
	_, err = fs.f.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("unable to get the beginning of file: %w", err)
	}

	err = json.NewEncoder(fs.f).Encode(&fs.urls)
	if err != nil {
		return fmt.Errorf("unable to encode data into the file: %w", err)
	}
	return nil
}

func NewFileStorage(filename string) (Storage, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, fmt.Errorf("unable to open file %s: %w", filename, err)
	}
	var short ShortURL
	urls := []*ShortURL{&short}
	if err := json.NewDecoder(file).Decode(&urls); err != nil && err != io.EOF {
		return nil, fmt.Errorf("unable to decode contents of file %s: %w", filename, err)
	}

	return &FileStorage{
		MemStorage: &MemStorage{urls: urls, userURLs: make(map[string][]*ShortURL)},
		f:          file,
	}, nil
}

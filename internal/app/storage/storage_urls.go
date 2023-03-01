package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"io"
	"os"
)

var ErrNotFound = errors.New("not found")
var ErrExistsURL = errors.New("url already exists")

type Storage interface {
	AddShortURL(userID string, shortURL *ShortURL) error
	GetShortURL(short string) (*ShortURL, error)
	GetAllShortURL(userID string) ([]*ShortURL, error)
	GetExistURL(originalURL string) (string, error)
}

type ShortURL struct {
	ID          string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"-"`
}

type MemStorage struct {
	urls []*ShortURL
}

func (ms *MemStorage) AddShortURL(userID string, s *ShortURL) error {
	s.UserID = userID
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

func (ms *MemStorage) GetAllShortURL(userID string) ([]*ShortURL, error) {
	var userShorts []*ShortURL
	for _, el := range ms.urls {
		if el.UserID == userID {
			userShorts = append(userShorts, el)
		}
	}
	return userShorts, nil
}

func (ms *MemStorage) GetExistURL(originalURL string) (string, error) {
	return "", nil
}

func NewMemoryStorage() Storage {
	var short = &ShortURL{"google", "https://www.google.com/", "12345"}
	return &MemStorage{
		urls: []*ShortURL{short},
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
		MemStorage: &MemStorage{urls: urls},
		f:          file,
	}, nil
}

type DBStorage struct {
	db *sql.DB
}

func (dbs *DBStorage) AddShortURL(userID string, s *ShortURL) error {
	_, err := dbs.db.Exec("INSERT INTO shorts (shortID, originalURL, userID) VALUES ($1, $2, $3)", s.ID, s.OriginalURL, userID)
	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			return ErrExistsURL
		}
		return err
	}
	return nil
}

func (dbs *DBStorage) GetShortURL(id string) (*ShortURL, error) {
	var short ShortURL
	row := dbs.db.QueryRow("SELECT shortID, originalURL FROM shorts WHERE shortID = $1", id)
	err := row.Scan(&short.ID, &short.OriginalURL)
	if err != nil {
		return nil, ErrNotFound
	}
	return &short, nil
}

func (dbs *DBStorage) GetAllShortURL(userID string) ([]*ShortURL, error) {
	rows, err := dbs.db.Query("SELECT shortID, originalURL FROM shorts WHERE userID = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	urls := make([]*ShortURL, 0)
	for rows.Next() {
		var short ShortURL
		err = rows.Scan(&short.ID, &short.OriginalURL)
		if err != nil {
			return nil, err
		}

		urls = append(urls, &short)

	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return urls, nil
}

func (dbs *DBStorage) GetExistURL(originalURL string) (string, error) {
	var short string
	row := dbs.db.QueryRow("SELECT shortID FROM shorts WHERE originalURL = $1", originalURL)
	err := row.Scan(&short)
	if err != nil {
		return "", err
	}
	return short, nil
}

func NewDB(dsn string) (Storage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to open sql connection: %w", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS shorts(shortID text PRIMARY KEY, originalURL text UNIQUE, userID text)")
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &DBStorage{
		db: db,
	}, nil
}

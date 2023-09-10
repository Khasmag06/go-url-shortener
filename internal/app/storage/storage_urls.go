package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Khasmag06/go-url-shortener/internal/app/models"
	"github.com/jackc/pgerrcode"
	"io"
	"os"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Общие ошибки
var (
	ErrNotFound     = errors.New("not found")
	ErrExistsURL    = errors.New("url already exists")
	ErrNotAvailable = errors.New("url removed")
)

// Storage интерфейс для работы с хранилищем.
type Storage interface {
	AddShortURL(userID string, shortURL *ShortURL) error
	GetShortURL(short string) (*ShortURL, error)
	GetAllShortURL(userID string) ([]*ShortURL, error)
	GetExistURL(originalURL string) (string, error)
	DeleteShortURL(userID, shortID string) error
	GetShortAndUserCount(stats *models.InternalStats) error
}

// ShortURL описание модели короткой ссылки.
type ShortURL struct {
	ID          string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"-"`
	IsDeleted   bool   `json:"-"`
}

// MemStorage структура, реализующая интерфейс Storage, для хранения данных в памяти.
type MemStorage struct {
	urls []*ShortURL
}

// AddShortURL добавляет короткую ссылку в память.
func (ms *MemStorage) AddShortURL(userID string, s *ShortURL) error {
	s.UserID = userID
	for _, el := range ms.urls {
		if el.OriginalURL == s.OriginalURL {
			return ErrExistsURL
		}
	}
	ms.urls = append(ms.urls, s)
	return nil

}

// GetShortURL возвращает оригинальную ссылку из памяти.
func (ms *MemStorage) GetShortURL(id string) (*ShortURL, error) {
	for _, el := range ms.urls {
		if el.ID == id {
			if !el.IsDeleted {
				return el, nil
			} else {
				return nil, ErrNotAvailable
			}
		}
	}

	return nil, ErrNotFound
}

// GetAllShortURL возвращает все ссылки из памяти.
func (ms *MemStorage) GetAllShortURL(userID string) ([]*ShortURL, error) {
	var userShorts []*ShortURL
	for _, el := range ms.urls {
		if el.UserID == userID {
			userShorts = append(userShorts, el)
		}
	}
	return userShorts, nil
}

// GetExistURL проверяет существует ли ссылка в памяти.
func (ms *MemStorage) GetExistURL(originalURL string) (string, error) {
	for _, el := range ms.urls {
		if el.OriginalURL == originalURL {
			return el.ID, nil
		}
	}
	return "", nil
}

// DeleteShortURL удаляет ссылку из памяти.
func (ms *MemStorage) DeleteShortURL(userID, shortURL string) error {
	for _, el := range ms.urls {
		if el.UserID == userID && el.ID == shortURL {
			el.IsDeleted = true
			return nil
		}
	}
	return nil
}

// GetShortAndUserCount получает количество коротких ссылок и пользователей из памяти.
func (ms *MemStorage) GetShortAndUserCount(stats *models.InternalStats) error {
	users := make(map[string]struct{})
	for _, el := range ms.urls {
		stats.ShortsCount++
		if _, ok := users[el.UserID]; !ok {
			users[el.UserID] = struct{}{}
			stats.UsersCount++
		}
	}
	return nil

}

// NewMemoryStorage конструктор для MemoryStorage.
func NewMemoryStorage() Storage {
	var short = &ShortURL{"google", "https://www.google.com/", "12345", false}
	return &MemStorage{
		urls: []*ShortURL{short},
	}
}

// FileStorage структура, реализующая интерфейс Storage для хранения данных в файле.
type FileStorage struct {
	*MemStorage
	f *os.File
}

// AddShortURL добавляет короткую ссылку в файл.
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

// NewFileStorage конструктор для FileStorage.
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

// DBStorage структура, реализующая интерфейс Storage, для хранения данных в базе данных.
type DBStorage struct {
	db *sql.DB
}

// AddShortURL добавляет короткую ссылку в бд.
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

// GetShortURL возвращает оригинальную ссылку из бд.
func (dbs *DBStorage) GetShortURL(id string) (*ShortURL, error) {
	var short ShortURL
	row := dbs.db.QueryRow("SELECT shortID, originalURL, is_deleted FROM shorts WHERE shortID = $1", id)
	err := row.Scan(&short.ID, &short.OriginalURL, &short.IsDeleted)
	if err != nil {
		return nil, ErrNotFound
	}
	if short.IsDeleted {
		return nil, ErrNotAvailable
	}
	return &short, nil
}

// GetAllShortURL возвращает все ссылки из бд.
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

// GetExistURL проверяет существует ли ссылка в бд.
func (dbs *DBStorage) GetExistURL(originalURL string) (string, error) {
	var short string
	row := dbs.db.QueryRow("SELECT shortID FROM shorts WHERE originalURL = $1", originalURL)
	err := row.Scan(&short)
	if err != nil {
		return "", err
	}
	return short, nil
}

// DeleteShortURL удаляет ссылку из бд.
func (dbs *DBStorage) DeleteShortURL(userID, shortID string) error {
	_, err := dbs.db.Exec("UPDATE shorts SET is_deleted=true WHERE userID = $1 AND shortID = $2", userID, shortID)
	if err != nil {
		return fmt.Errorf("unable to delete URL with %s: %w", shortID, err)
	}
	return nil
}

// GetShortAndUserCount получает количество коротких ссылок и пользователей из бд.
func (dbs *DBStorage) GetShortAndUserCount(stats *models.InternalStats) error {
	row := dbs.db.QueryRow("SELECT COUNT(shortID), COUNT(DISTINCT userID) FROM shorts")
	err := row.Scan(&stats.ShortsCount, &stats.UsersCount)
	if err != nil {
		return err
	}
	return nil
}

// NewDB конструктор для DBStorage
func NewDB(dsn string) (Storage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to open sql connection: %w", err)
	}

	query := `CREATE TABLE IF NOT EXISTS shorts(shortID text PRIMARY KEY, 
										 originalURL text UNIQUE,
										 userID text,
                                         is_deleted BOOLEAN NOT NULL DEFAULT false)`

	_, err = db.Exec(query)
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

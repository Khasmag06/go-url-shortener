package models

// InternalStats описания модели ответа количества сокращенных ссылок и пользователей.
type InternalStats struct {
	ShortsCount int `json:"urls"`
	UsersCount  int `json:"users"`
}

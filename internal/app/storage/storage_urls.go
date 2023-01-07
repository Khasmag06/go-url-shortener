package storage

type Storage struct {
	urls map[string]string
}

func (s *Storage) Put(id, url string) {
	s.urls[id] = url
}

func (s *Storage) Get(id string) string {
	url := s.urls[id]
	return url
}

var Urls = Storage{make(map[string]string)}

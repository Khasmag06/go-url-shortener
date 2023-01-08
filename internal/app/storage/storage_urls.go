package storage

type Storage struct {
	urls map[string]string
}

func (s *Storage) Put(id, url string) {
	s.urls[id] = url
}

func (s *Storage) Get(id string) string {
	url, ok := s.urls[id]
	if !ok {
		return ""
	}
	return url

}

var Urls = Storage{make(map[string]string)}

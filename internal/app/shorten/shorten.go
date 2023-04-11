package shorten

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var length = 6

// URLShorten создает случайную строку длинною 6 латинский букв.
func URLShorten() string {
	var url strings.Builder
	for i := 0; i < length; i++ {
		url.WriteByte(alphabet[rand.Intn(len(alphabet))])
	}
	return url.String()
}

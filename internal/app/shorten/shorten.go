package shorten

import (
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var length = 6

func URLShorten() string {
	var url strings.Builder
	for i := 0; i < length; i++ {
		url.WriteByte(alphabet[rand.Intn(len(alphabet))])
	}
	return "/" + url.String()
}

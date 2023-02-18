package middleware

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
	"io"
	"net/http"
)

var key = sha256.Sum256([]byte("Secret key"))

//var session = map[string]struct{}{"12345": {}}

func CreateAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err == http.ErrNoCookie {
			userID := uuid.NewString()
			//session[userID] = struct{}{}
			cookie = &http.Cookie{
				Name:   "token",
				Value:  encrypt(userID),
				Path:   "/",
				MaxAge: 300,
			}
		}

		//if _, ok := session[decrypt(cookie.Value)]; !ok {
		//	userID := uuid.NewString()
		//	session[userID] = struct{}{}
		//	cookie = &http.Cookie{
		//		Name:   "token",
		//		Value:  encrypt(userID),
		//		Path:   "/",
		//		MaxAge: 300,
		//	}
		//}
		http.SetCookie(w, cookie)
		ctx := context.WithValue(r.Context(), "userID", decrypt(cookie.Value))
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func createGcm() cipher.AEAD {
	aesBlock, err := aes.NewCipher(key[:])
	if err != nil {
		panic(err)
	}
	aesGcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		panic(err)
	}
	return aesGcm
}

func encrypt(data string) string {
	aesGcm := createGcm()
	nonce := make([]byte, aesGcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	ciphertext := aesGcm.Seal(nonce, nonce, []byte(data), nil)
	return hex.EncodeToString(ciphertext)

}

func decrypt(encrypt string) string {
	aesGcm := createGcm()
	data, err := hex.DecodeString(encrypt)
	if err != nil {
		panic(err)
	}
	nonceSize := aesGcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesGcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err)
	}
	return string(plaintext)

}

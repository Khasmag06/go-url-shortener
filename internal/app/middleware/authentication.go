package middleware

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

var aesGcm cipher.AEAD

func init() {
	var err error
	aesGcm, err = createGcm()
	if err != nil {
		log.Fatal(err)
	}
}

type userIDKeyType string

const UserIDKey userIDKeyType = "userID"

var key = sha256.Sum256([]byte("Secret key"))

func CreateAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err == http.ErrNoCookie {
			userID := uuid.NewString()
			cookie = &http.Cookie{
				Name:   "token",
				Value:  encrypt(userID),
				Path:   "/",
				MaxAge: 300,
			}
		}
		http.SetCookie(w, cookie)
		cookieDecrypt, err := decrypt(cookie.Value)
		if err != nil {
			log.Fatal(err)
		}
		ctx := context.WithValue(r.Context(), UserIDKey, cookieDecrypt)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func createGcm() (cipher.AEAD, error) {
	aesBlock, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("unable to initialize new cipher: %w", err)
	}
	aesGcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize new GCM: %w", err)
	}
	return aesGcm, nil
}

func encrypt(data string) string {
	nonce := make([]byte, aesGcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	ciphertext := aesGcm.Seal(nonce, nonce, []byte(data), nil)
	return hex.EncodeToString(ciphertext)

}

func decrypt(encrypt string) (string, error) {
	data, err := hex.DecodeString(encrypt)
	if err != nil {
		return "", fmt.Errorf("unable to decode string: %w", err)
	}
	nonceSize := aesGcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesGcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("unable to get encrypted data: %w", err)
	}
	return string(plaintext), nil

}

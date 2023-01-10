package handlers_test

import (
	"fmt"
	"github.com/Khasmag06/go-url-shortener/internal/app/handlers"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHandler(t *testing.T) {
	type want struct {
		code     int
		location string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code:     http.StatusTemporaryRedirect,
				location: "https://www.google.com/",
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/google", nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handlers.GetHandler)
			h.ServeHTTP(w, request)
			response := w.Result()
			defer response.Body.Close()
			fmt.Println(response.StatusCode)
			fmt.Println(response.Header.Get("Location"))
			assert.Equal(t, response.StatusCode, tt.want.code)
			assert.Equal(t, response.Header.Get("Location"), tt.want.location)

		})
	}
}

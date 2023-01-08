package handlers_test

import (
	"bytes"
	"github.com/Khasmag06/go-url-shortener/internal/app/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
)

func TestPostHandler(t *testing.T) {
	type want struct {
		code        int
		contentType string
		response    string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code:        http.StatusCreated,
				contentType: "text/plain; charset=utf-8",
				response:    `^http://localhost:8080/([a-zA-Z]{6})$`,
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := url.Values{}
			long := "https://www.google.com/"
			data.Set("url", long)
			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(data.Encode()))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handlers.PostHandler)
			h.ServeHTTP(w, request)
			response := w.Result()

			assert.Equal(t, response.StatusCode, tt.want.code)
			assert.Equal(t, response.Header.Get("Content-Type"), tt.want.contentType)
			body, err := io.ReadAll(response.Body)
			require.NoError(t, err)
			err = response.Body.Close()
			require.NoError(t, err)
			assert.Regexp(t, regexp.MustCompile(tt.want.response), string(body))

		})
	}
}

package main

import (
	"github.com/Khasmag06/go-url-shortener/config"
	"github.com/Khasmag06/go-url-shortener/internal/app/handlers"
	"github.com/Khasmag06/go-url-shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	repo := storage.NewMemoryStorage()
	s := handlers.NewService(*cfg, repo)
	r := NewRouter(s)
	ts := httptest.NewServer(r)
	defer ts.Close()

	statusCode, _ := testRequest(t, ts, "GET", "/google")
	assert.Equal(t, http.StatusOK, statusCode)
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp.StatusCode, string(respBody)
}

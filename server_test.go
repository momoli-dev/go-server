package server_test

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/momoli-dev/go-server"
	"github.com/stretchr/testify/assert"
)

func TestNewServer_OK(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	params := &server.ServerParams{
		Addr:    ":8081",
		Handler: handler,
	}

	srv := server.NewServer(params)

	go func() {
		err := srv.Start()
		assert.NoError(t, err)
	}()
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:8081/test")
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	expectedBody := "Hello, World!"
	assert.Equal(t, expectedBody, string(body))

	err = srv.Shutdown()
	assert.NoError(t, err)

	_, err = http.Get("http://localhost:8081/test")
	assert.Error(t, err)
}

func TestNewServer_InvalidAddr(t *testing.T) {
	handler := http.NewServeMux()
	params := &server.ServerParams{
		Addr:    "invalid-addr",
		Handler: handler,
	}

	srv := server.NewServer(params)

	err := srv.Start()
	assert.Error(t, err)
}

func TestNewServer_NilHandler(t *testing.T) {
	params := &server.ServerParams{
		Addr:    ":8082",
		Handler: nil,
	}

	assert.Panics(t, func() {
		server.NewServer(params)
	})
}

func TestServer_MustStart_Panic(t *testing.T) {
	handler := http.NewServeMux()
	params := &server.ServerParams{
		Addr:    "invalid-addr",
		Handler: handler,
	}

	srv := server.NewServer(params)

	assert.Panics(t, func() {
		srv.MustStart()
	})
}

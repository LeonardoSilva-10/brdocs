package httpclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	t.Run("uses default timeout", func(t *testing.T) {
		c := NewClient(0)
		assert.Equal(t, DefaultTimeout, c.Timeout)
	})

	t.Run("uses custom timeout", func(t *testing.T) {
		c := NewClient(5 * time.Second)
		assert.Equal(t, 5*time.Second, c.Timeout)
	})
}

func TestReadBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	require.NoError(t, err)

	body, err := ReadBody(resp)
	require.NoError(t, err)
	assert.Equal(t, "hello world", string(body))
}

func TestCheckStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		expected []int
		wantErr  bool
	}{
		{"ok single", http.StatusOK, []int{http.StatusOK}, false},
		{"ok multiple", http.StatusCreated, []int{http.StatusOK, http.StatusCreated}, false},
		{"error", http.StatusNotFound, []int{http.StatusOK}, true},
		{"empty expected", http.StatusOK, []int{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{StatusCode: tt.status}
			err := CheckStatus(resp, tt.expected...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDoWithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	t.Run("successful request", func(t *testing.T) {
		client := NewClient(5 * time.Second)
		req, err := http.NewRequest(http.MethodGet, server.URL, nil)
		require.NoError(t, err)

		resp, err := DoWithContext(context.Background(), client, req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	})

	t.Run("context cancellation", func(t *testing.T) {
		client := NewClient(5 * time.Second)
		req, err := http.NewRequest(http.MethodGet, server.URL, nil)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		_, err = DoWithContext(ctx, client, req)
		assert.Error(t, err)
	})
}

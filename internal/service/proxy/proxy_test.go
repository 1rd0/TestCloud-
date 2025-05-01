package proxy

import (
	"net/http"
	"net/http/httptest"
	"net/url" // Добавьте этот импорт
	"testing"

	"github.com/1rd0/TestCloud-/internal/service/backend"
)

func TestHandler_ServeHTTP(t *testing.T) {
	// Создаем тестовый сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer ts.Close() // Не забываем закрыть сервер

	// Преобразуем строку URL в *url.URL
	backendURL, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatalf("Failed to parse test server URL: %v", err)
	}

	b := backend.New(backendURL) // Теперь передаем *url.URL, а не строку

	h := New(func() (*backend.Backend, error) {
		return b, nil
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rw.Code)
	}
}

func BenchmarkHandler_ServeHTTP(b *testing.B) {
	// Создаем тестовый сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close() // Не забываем закрыть сервер

	// Преобразуем строку URL в *url.URL
	backendURL, err := url.Parse(ts.URL)
	if err != nil {
		b.Fatalf("Failed to parse test server URL: %v", err)
	}

	bck := backend.New(backendURL) // Теперь передаем *url.URL, а не строку

	h := New(func() (*backend.Backend, error) {
		return bck, nil
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	for i := 0; i < b.N; i++ {
		rw := httptest.NewRecorder()
		h.ServeHTTP(rw, req)
	}
}

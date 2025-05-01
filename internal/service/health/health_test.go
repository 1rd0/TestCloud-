package health

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url" // Добавьте этот импорт
	"testing"
	"time"

	"github.com/1rd0/TestCloud-/internal/service/backend"
	"go.uber.org/zap"
)

func TestStartHealthCheck(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Преобразуем строку URL в *url.URL
	backendURL, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatalf("Failed to parse test server URL: %v", err)
	}

	b := backend.New(backendURL) // Теперь передаем *url.URL, а не строку
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log, _ := zap.NewDevelopment()
	Start(ctx, []*backend.Backend{b}, 100*time.Millisecond, 100*time.Millisecond, log)

	time.Sleep(300 * time.Millisecond)

	if !b.IsAlive() {
		t.Error("backend expected to be alive")
	}
}

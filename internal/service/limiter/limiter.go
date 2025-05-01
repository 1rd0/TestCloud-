package limiter

import (
	"context"

	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"net"
	"net/http"

	"sync"
)

type ConfigRow struct {
	ID           string
	Capacity     int
	RefillPerSec float64
}

type Limiter struct {
	db      *pgxpool.Pool
	buckets map[string]*Bucket
	mu      sync.Mutex
}

func NewLimiter(ctx context.Context, pool *pgxpool.Pool) (*Limiter, error) {
	return &Limiter{
		db:      pool,
		buckets: make(map[string]*Bucket),
	}, nil
}

func (l *Limiter) getConfig(ctx context.Context, clientID string) (*ConfigRow, error) {
	var (
		capacity   int
		ratePerSec int
		cfgID      string
	)

	err := l.db.QueryRow(ctx,
		`SELECT id, capacity, rate_per_sec FROM clients WHERE id=$1`, clientID,
	).Scan(&cfgID, &capacity, &ratePerSec)
	if err != nil {

		err = l.db.QueryRow(ctx,
			`SELECT id, capacity, rate_per_sec FROM clients WHERE id='default'`,
		).Scan(&cfgID, &capacity, &ratePerSec)
		if err != nil {
			return nil, err
		}
	}

	return &ConfigRow{
		ID:           cfgID,
		Capacity:     capacity,
		RefillPerSec: float64(ratePerSec),
	}, nil
}

func (l *Limiter) getBucket(ctx context.Context, clientID string) (*Bucket, error) {
	l.mu.Lock()
	bucket, ok := l.buckets[clientID]
	l.mu.Unlock()
	if ok {
		return bucket, nil
	}

	cfg, err := l.getConfig(ctx, clientID)
	if err != nil {
		return nil, err
	}
	bucket = NewBucket(cfg.Capacity, cfg.RefillPerSec)

	l.mu.Lock()
	l.buckets[clientID] = bucket
	l.mu.Unlock()
	return bucket, nil
}

func (l *Limiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		clientID := r.Header.Get("X-API-Key")
		if clientID == "" {
			host, _, _ := net.SplitHostPort(r.RemoteAddr)
			clientID = host
		}

		bucket, err := l.getBucket(r.Context(), clientID)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		if err := bucket.TryTake(); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    429,
				"message": "Rate limit exceeded",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

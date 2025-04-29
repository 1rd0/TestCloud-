package middleware

import (
	"encoding/json"
	"net/http"

	"your/app/internal/limiter"
)

func RateLimit(l *limiter.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := limiter.ClientID(r.RemoteAddr)
			if apiKey := r.Header.Get("X-Api-Key"); apiKey != "" {
				id = apiKey
			}

			if !l.Allow(id) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"code":    429,
					"message": "Rate limit exceeded",
				})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

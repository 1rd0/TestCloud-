package middleware

import (
	"context"
	mylog "github.com/1rd0/TestCloud-/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// small wrapper чтобы поймать статус/байты
type recorder struct {
	http.ResponseWriter
	status int
	size   int
}

func (r *recorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *recorder) Write(b []byte) (int, error) {
	n, err := r.ResponseWriter.Write(b)
	r.size += n
	return n, err
}

func AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. request id
		reqID := uuid.NewString()
		ctx := context.WithValue(r.Context(), mylog.RequestId, reqID)

		rec := &recorder{ResponseWriter: w, status: http.StatusOK}
		start := time.Now()

		// 2. передать дальше по цепочке
		next.ServeHTTP(rec, r.WithContext(ctx))

		// 3. лог после обработки
		l := mylog.GetLoggerFromCtx(ctx)
		if l == nil { // safety
			return
		}

		backend := "-"
		if v := r.Context().Value("backend"); v != nil {
			backend = v.(string)
		}

		l.Info(ctx, "request completed",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.String("backend", backend),
			zap.Int("status", rec.status),
			zap.Int("bytes", rec.size),
			zap.Duration("latency", time.Since(start)),
			zap.String("client_ip", r.RemoteAddr),
		)
	})
}

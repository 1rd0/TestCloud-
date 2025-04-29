package health

import (
	"context"
	"net"
	"time"

	"go.uber.org/zap"
)

func Start(ctx context.Context, addrs []string, interval time.Duration, log *zap.Logger) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				for _, a := range addrs {
					conn, err := net.DialTimeout("tcp", a, 2*time.Second)
					if err != nil {
						log.Warn("backend down", zap.String("addr", a), zap.Error(err))
						continue
					}
					 
					_ = conn.Close()
				}
			}
		}
	}()
}

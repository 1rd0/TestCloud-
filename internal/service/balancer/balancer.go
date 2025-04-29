package balancer

import (
	"github.com/1rd0/TestCloud-/internal/service/backend"
)

type Balancer interface {
	// Next возвращает URL следующего живого бэкенда или ошибку,
	// если все бэкенды помечены как мертвые.
	Next() (*backend.Backend, error)
}

package balancer

import "net/url"

type Balancer interface {
	Next() (*url.URL, error)
}

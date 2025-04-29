package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	TotalRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "lb_total_requests",
			Help: "total HTTP requests processed by the load balancer",
		})
	BackendUp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "lb_backend_up",
			Help: "backend health-state (1 = up, 0 = down)",
		},
		[]string{"backend"},
	)
)

func init() {
	prometheus.MustRegister(TotalRequests)
	prometheus.MustRegister(BackendUp)
}

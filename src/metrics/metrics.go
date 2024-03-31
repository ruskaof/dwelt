package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// counter of each http request by method and path
	httpRequestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "path"})

	httpRequestDurationVec = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "http_request_duration_seconds",
		Help:       "Duration of HTTP requests",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001, 1: 0},
	}, []string{"method", "path"})

	incomingHttpBytes = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_incoming_bytes_total",
		Help: "Total number of outgoing bytes",
	}, []string{"method", "path"})

	incomingWebsocketBytes = promauto.NewCounter(prometheus.CounterOpts{
		Name: "websocket_incoming_bytes_total",
		Help: "Total number of incoming bytes",
	})

	// number of websocket connections
	websocketConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "websocket_connections",
		Help: "Number of websocket connections",
	})
)

// RecordHttpRequest records the duration of the request
func RecordHttpRequest(method, path string, duration float64, incomingBytes int) {
	httpRequestCounter.WithLabelValues(method, path).Inc()
	httpRequestDurationVec.WithLabelValues(method, path).Observe(duration)
	incomingHttpBytes.WithLabelValues(method, path).Add(float64(incomingBytes))
}

func IncrementWebsocketConnections() {
	websocketConnections.Inc()
}

func DecrementWebsocketConnections() {
	websocketConnections.Dec()
}

func IncrementIncomingWebsocketBytes(bytes int) {
	incomingWebsocketBytes.Add(float64(bytes))
}

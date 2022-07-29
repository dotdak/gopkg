package metrics

import (
	"context"
	"net/http"

	"github.com/dotdak/gopkg/env"
	"github.com/dotdak/gopkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewMetricServer(
	metricPort string,
) *MetricServer {
	port := env.EnvString("METRIC_PORT", metricPort)
	logger.LOG().Info("metric server", "address", port)
	return &MetricServer{
		Server: &http.Server{Addr: port, Handler: nil},
	}
}

type MetricServer struct {
	*http.Server
}

func (s *MetricServer) Start() error {
	s.Server.Handler = promhttp.Handler()
	return s.Server.ListenAndServe()
}

func (s *MetricServer) Stop() error {
	return s.Server.Shutdown(context.Background())
}

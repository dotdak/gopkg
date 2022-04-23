package metrics

import (
	"context"
	"net/http"

	"github.com/dotdak/gopkg/env"
	"github.com/dotdak/gopkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewMetricServer(
	ctx context.Context,
) *MetricServer {
	metricPort := env.EnvString("METRIC_PORT", ":7070")
	logger.LOG().Info("metric server", "address", metricPort)
	return &MetricServer{
		Server: &http.Server{Addr: metricPort, Handler: nil},
	}
}

type MetricServer struct {
	*http.Server
}

func (s *MetricServer) Start() error {
	s.Server.Handler = promhttp.Handler()
	return s.ListenAndServe()
}

func (s *MetricServer) Stop() error {
	return s.Shutdown(context.Background())
}

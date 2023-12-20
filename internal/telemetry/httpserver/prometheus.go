package httpserver

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/httpserver"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type prometheusBuilder struct {
	registry *prometheus.Registry
}

func (prometheusBuilder) Key() string {
	return "prometheus"
}

func (b prometheusBuilder) Apply(e *gin.Engine, cfg httpserver.Config) error {
	h := promhttp.HandlerFor(b.registry, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	})
	e.Any("/metrics", func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	})
	return nil
}

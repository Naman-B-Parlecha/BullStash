package main

import (
	"github.com/Naman-B-Parlecha/BullStash/monitoring/models"
	"github.com/Naman-B-Parlecha/BullStash/monitoring/routes"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var metrics *models.Metrics

func main() {
	reg := prometheus.NewRegistry()
	metrics = models.NewMetrics(reg)

	r := gin.Default()
	r.GET("/metrics", gin.WrapH(promhttp.HandlerFor(reg, promhttp.HandlerOpts{})))

	routes.SetupRoutes(r, metrics)

	r.Run(":8085")
}

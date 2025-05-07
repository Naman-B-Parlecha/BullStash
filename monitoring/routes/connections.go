package routes

import (
	"net/http"

	"github.com/Naman-B-Parlecha/BullStash/monitoring/models"
	"github.com/gin-gonic/gin"
)

func ConnectionRoutes(r *gin.RouterGroup, metrics *models.Metrics) {
	r.POST("/success", func(c *gin.Context) {
		var metric struct {
			DBType string `json:"dbtype"`
		}

		if err := c.ShouldBindJSON(&metric); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Bad payload"})
			return
		}

		metrics.ConnectionSuccessCount.WithLabelValues(metric.DBType).Inc()
		c.JSON(http.StatusOK, gin.H{"message": "Successfully updated monitoring"})
	})

	r.POST("/failure", func(c *gin.Context) {
		var metric struct {
			DBType string `json:"dbtype"`
		}

		if err := c.ShouldBindJSON(&metric); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Bad payload"})
			return
		}

		metrics.ConnectionFailureCount.WithLabelValues(metric.DBType).Inc()
		c.JSON(http.StatusOK, gin.H{"message": "Successfully updated monitoring"})
	})

	r.POST("/latency", func(c *gin.Context) {
		var metric struct {
			DBType   string  `json:"dbtype"`
			Duration float64 `json:"duration"`
		}

		if err := c.ShouldBindJSON(&metric); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Bad payload"})
			return
		}

		if metric.Duration < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Duration must be non-negative"})
			return
		}

		metrics.ConnectionLatency.WithLabelValues(metric.DBType).Observe(metric.Duration)
		c.JSON(http.StatusOK, gin.H{"message": "Successfully updated monitoring"})
	})
}

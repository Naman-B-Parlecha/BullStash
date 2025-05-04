package routes

import (
	"net/http"

	"github.com/Naman-B-Parlecha/BullStash/monitoring/models"
	"github.com/gin-gonic/gin"
)

func RestoreRoutes(r *gin.RouterGroup, metrics *models.Metrics) {
	r.POST("/success", func(c *gin.Context) {
		var metric struct {
			DBType  string `json:"dbtype"`
			Storage string `json:"storage"`
		}

		if err := c.ShouldBindJSON(&metric); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Bad payload"})
			return
		}

		metrics.RestoreSuccessCount.WithLabelValues(metric.DBType, metric.Storage).Inc()
		c.JSON(http.StatusOK, gin.H{"messaage": "Successfully updated monitoring"})
	})

	r.POST("/failure", func(c *gin.Context) {
		var metric struct {
			DBType  string `json:"dbtype"`
			Storage string `json:"storage"`
		}

		if err := c.ShouldBindJSON(&metric); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Bad payload"})
			return
		}

		metrics.RestoreFailureCount.WithLabelValues(metric.DBType, metric.Storage).Inc()
		c.JSON(http.StatusOK, gin.H{"messaage": "Successfully updated monitoring"})
	})

	r.POST("/duration", func(c *gin.Context) {
		var metric struct {
			DBType   string  `json:"dbtype"`
			Duration float64 `json:"duration"`
			Storage  string  `json:"storage"`
		}

		if err := c.ShouldBindJSON(&metric); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Bad payload"})
			return
		}

		metrics.RestoreDuration.WithLabelValues(metric.DBType, metric.Storage).Observe(metric.Duration)
		c.JSON(http.StatusOK, gin.H{"messaage": "Successfully updated monitoring"})
	})
}

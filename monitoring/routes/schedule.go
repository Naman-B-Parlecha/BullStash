package routes

import (
	"net/http"

	"github.com/Naman-B-Parlecha/BullStash/monitoring/models"
	"github.com/gin-gonic/gin"
)

func ScheduleRoutes(r *gin.RouterGroup, metrics *models.Metrics) {
	r.POST("/executioncount", func(c *gin.Context) {
		type metric struct {
			DBType string `json:"dbtype"`
		}

		if err := c.ShouldBindJSON(&metric{}); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		metrics.ScheduledBackupExecutionCount.WithLabelValues("dbtype").Inc()
		c.JSON(http.StatusOK, gin.H{"message": "Scheduled backup execution count incremented"})
	})

	r.POST("/success", func(c *gin.Context) {
		type metric struct {
			DBType string `json:"dbtype"`
		}

		if err := c.ShouldBindJSON(&metric{}); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		metrics.ScheduledBackupSuccessCount.WithLabelValues("dbtype").Inc()
		c.JSON(http.StatusOK, gin.H{"message": "Scheduled backup success count incremented"})
	})

	r.POST("/failure", func(c *gin.Context) {
		type metric struct {
			DBType string `json:"dbtype"`
		}

		if err := c.ShouldBindJSON(&metric{}); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		metrics.ScheduledBackupFailureCount.WithLabelValues("dbtype").Inc()
		c.JSON(http.StatusOK, gin.H{"message": "Scheduled backup failure count incremented"})
	})
}

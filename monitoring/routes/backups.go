package routes

import (
	"net/http"

	"github.com/Naman-B-Parlecha/BullStash/monitoring/models"
	"github.com/gin-gonic/gin"
)

func BackupRoutes(r *gin.RouterGroup, metrics *models.Metrics) {
	r.POST("/success", func(c *gin.Context) {
		var metric struct {
			DBType      string `json:"dbtype"`
			Backup_type string `json:"backup_type"`
			Storage     string `json:"storage"`
		}

		if err := c.ShouldBindJSON(&metric); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Bad payload"})
			return
		}

		metrics.BackupSuccessCount.WithLabelValues(metric.DBType, metric.Backup_type, metric.Storage).Inc()
		c.JSON(http.StatusOK, gin.H{"messaage": "Successfully updated monitoring"})
	})

	r.POST("/failure", func(c *gin.Context) {
		var metric struct {
			DBType      string `json:"dbtype"`
			Backup_type string `json:"backup_type"`
			Storage     string `json:"storage"`
		}

		if err := c.ShouldBindJSON(&metric); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Bad payload"})
			return
		}

		metrics.BackupFailureCount.WithLabelValues(metric.DBType, metric.Backup_type, metric.Storage).Inc()
		c.JSON(http.StatusOK, gin.H{"messaage": "Successfully updated monitoring"})
	})

	r.POST("/duration", func(c *gin.Context) {
		var metric struct {
			DBType     string  `json:"dbtype"`
			BackupType string  `json:"backup_type"`
			Storage    string  `json:"storage"`
			Duration   float64 `json:"duration"`
		}

		if err := c.ShouldBindJSON(&metric); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Bad payload"})
			return
		}

		metrics.BackupDuration.WithLabelValues(metric.DBType, metric.BackupType, metric.Storage).Observe(metric.Duration)
		c.JSON(http.StatusOK, gin.H{"messaage": "Successfully updated monitoring"})
	})

	r.POST("/size", func(c *gin.Context) {
		var metric struct {
			DBType     string  `json:"dbtype"`
			Size       float64 `json:"size"`
			BackupType string  `json:"backup_type"`
			Storage    string  `json:"storage"`
		}

		if err := c.ShouldBindJSON(&metric); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Bad payload"})
			return
		}

		metrics.BackupSize.WithLabelValues(metric.DBType, metric.BackupType, metric.Storage).Observe(metric.Size)
		c.JSON(http.StatusOK, gin.H{"messaage": "Successfully updated monitoring"})
	})
}

package routes

import (
	"github.com/Naman-B-Parlecha/BullStash/monitoring/models"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, metrics *models.Metrics) {
	backupRoute := r.Group("/backups")
	{
		BackupRoutes(backupRoute, metrics)
	}

	restoreRoute := r.Group("/restore")
	{
		RestoreRoutes(restoreRoute, metrics)
	}
	connectionRoute := r.Group("/connections")
	{
		ConnectionRoutes(connectionRoute, metrics)
	}
	scheduleRoute := r.Group("/schedules")
	{
		ScheduleRoutes(scheduleRoute, metrics)
	}
	compressRoute := r.Group("/compress")
	{
		CompressRoutes(compressRoute, metrics)
	}
}

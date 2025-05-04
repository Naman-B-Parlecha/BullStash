package models

import (
	"math"
	"slices"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	BackupDuration                *prometheus.HistogramVec
	BackupSize                    *prometheus.HistogramVec
	BackupSuccessCount            *prometheus.CounterVec
	BackupFailureCount            *prometheus.CounterVec
	RestoreDuration               *prometheus.HistogramVec
	RestoreSuccessCount           *prometheus.CounterVec
	RestoreFailureCount           *prometheus.CounterVec
	ConnectionLatency             *prometheus.HistogramVec
	ConnectionFailureCount        *prometheus.CounterVec
	ConnectionSuccessCount        *prometheus.CounterVec
	ScheduledBackupExecutionCount *prometheus.CounterVec
	ScheduledBackupSuccessCount   *prometheus.CounterVec
	ScheduledBackupFailureCount   *prometheus.CounterVec
	CompressionDuration           *prometheus.HistogramVec
	CompressionRatio              *prometheus.GaugeVec
	CPUUsage                      *prometheus.GaugeVec
	MemoryUsage                   *prometheus.GaugeVec
	DiskUsage                     *prometheus.GaugeVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		BackupDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "bullstash_backup_duration_seconds",
			Help:      "Duration of the backup process in seconds",
			Namespace: "myapp",
			Buckets:   []float64{0.1, 0.5, 1, 5, 10, 30, 60},
		}, []string{"dbtype", "backup_type", "storage"}),
		BackupSize: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "bullstash_backup_size_bytes",
			Help:      "Size of the backup in bytes",
			Namespace: "myapp",
			Buckets:   []float64{0, 10, 50, 100, 1000, 5000, 10000, 50000, 100000, 500000, 1000000},
		}, []string{"dbtype", "backup_type", "storage"}),
		BackupSuccessCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "bullstash_backup_success_count",
			Help:      "Number of successful backups",
			Namespace: "myapp",
		}, []string{"dbtype", "backup_type", "storage"}),
		BackupFailureCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "bullstash_backup_failure_count",
			Help:      "Number of failed backups",
			Namespace: "myapp",
		}, []string{"dbtype", "backup_type", "storage"}),

		RestoreDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "bullstash_restore_duration_seconds",
			Help:      "Duration of the restore process in seconds",
			Namespace: "myapp",
			Buckets:   []float64{0.1, 0.5, 1, 5, 10, 30, 60},
		}, []string{"dbtype", "storage"}),
		RestoreSuccessCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "bullstash_restore_success_count",
			Help:      "Number of successful restores",
			Namespace: "myapp",
		}, []string{"dbtype", "storage"}),
		RestoreFailureCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "bullstash_restore_failure_count",
			Help:      "Number of failed restores",
			Namespace: "myapp",
		}, []string{"dbtype", "storage"}),

		ConnectionLatency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "bullstash_connection_latency_seconds",
			Help:      "Latency of database connection attempts in seconds",
			Namespace: "myapp",
			Buckets:   []float64{0.01, 0.05, 0.1, 0.5, 1, 5},
		}, []string{"dbtype"}),
		ConnectionFailureCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "bullstash_connection_failure_count",
			Help:      "Number of failed database connection attempts",
			Namespace: "myapp",
		}, []string{"dbtype"}),
		ConnectionSuccessCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "bullstash_connection_success_count",
			Help:      "Number of successful database connection attempts",
			Namespace: "myapp",
		}, []string{"dbtype"}),

		ScheduledBackupExecutionCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "bullstash_scheduled_backup_execution_count",
			Help:      "Number of scheduled backup executions",
			Namespace: "myapp",
		}, []string{"dbtype"}),
		ScheduledBackupSuccessCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "bullstash_scheduled_backup_success_count",
			Help:      "Number of successful scheduled backups",
			Namespace: "myapp",
		}, []string{"dbtype"}),
		ScheduledBackupFailureCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "bullstash_scheduled_backup_failure_count",
			Help:      "Number of failed scheduled backups",
			Namespace: "myapp",
		}, []string{"dbtype"}),

		CompressionDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "bullstash_compression_duration_seconds",
			Help:      "Duration of the compression process in seconds",
			Namespace: "myapp",
			Buckets:   []float64{0.01, 0.05, 0.1, 0.5, 1, 5},
		}, []string{"dbtype", "storage"}),
		CompressionRatio: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "bullstash_compression_ratio",
			Help:      "Compression ratio achieved during backup (compressed size / original size)",
			Namespace: "myapp",
		}, []string{"dbtype", "storage"}),

		CPUUsage: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "bullstash_cpu_usage_percent",
			Help:      "CPU usage of the application in percent",
			Namespace: "myapp",
		}, []string{"dbtype", "storage"}),
		MemoryUsage: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "bullstash_memory_usage_bytes",
			Help:      "Memory usage of the application in bytes",
			Namespace: "myapp",
		}, []string{"dbtype", "storage"}),
		DiskUsage: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "bullstash_disk_usage_bytes",
			Help:      "Disk usage of the application in bytes",
			Namespace: "myapp",
		}, []string{"dbtype", "storage"}),
	}

	reg.MustRegister(
		m.BackupDuration, m.BackupFailureCount, m.BackupSuccessCount, m.BackupSize,
		m.RestoreDuration, m.RestoreFailureCount, m.RestoreSuccessCount,
		m.ConnectionLatency, m.ConnectionFailureCount, m.ConnectionSuccessCount,
		m.ScheduledBackupExecutionCount, m.ScheduledBackupFailureCount, m.ScheduledBackupSuccessCount,
		m.CompressionDuration, m.CompressionRatio,
		m.CPUUsage, m.MemoryUsage, m.DiskUsage,
	)

	return m
}

func (m *Metrics) AddDummyData() {
	m.BackupSuccessCount.WithLabelValues("postgres", "full", "local").Add(3)
	m.BackupFailureCount.WithLabelValues("postgres", "full", "local").Add(2)
	m.BackupSuccessCount.WithLabelValues("mongodb", "full", "local").Add(4)
	m.BackupFailureCount.WithLabelValues("mongodb", "full", "local").Add(3)
	m.BackupSuccessCount.WithLabelValues("mysql", "full", "local").Add(1)
	m.BackupFailureCount.WithLabelValues("mysql", "full", "local").Add(2)

	durationObservations := []float64{1.5, 3, 4.5, 1.5, 7.5,
		2.1, 3.4, 5.7, 8.9, 10.2,
		1.1, 2.3, 4.2, 6.7, 9.8,
		11.5, 7.8, 3.3, 4.4, 5.5,
		6.6, 8.2, 9.1, 10.7, 2.8,
		5.9, 8.3, 11.2, 3.7, 1.9,
		4.8, 7.2, 9.3, 11.8, 2.2,
		5.1, 8.7, 10.5, 3.8, 1.3,
		6.4, 9.5, 11.9, 2.7, 4.1,
		7.3, 9.9, 1.7, 5.2, 8.1,
		10.9, 12.0, 3.9, 6.8, 9.4}

	sizeObservations := []float64{12345678, 890123, 45678901, 2345678, 123456789,
		567890, 3456789, 98765432, 876543, 23456789,
		1234567, 87654321, 345678, 9876543, 234567,
		56789012, 3456789012, 45678, 987654, 2345,
		789012345, 67890, 1234, 56789, 9876,
		3456789012, 456789, 23456, 8901234, 567,
		12345, 678901, 2345678901, 34567, 890,
		45678, 1234567890, 23456789, 34567890, 4567890,
		567891234, 78901, 234567890, 345678, 45678901,
		5678, 67890123, 789012, 89012345, 9012345}

	x := sizeObservations
	slices.Reverse(sizeObservations)

	go func() {
		for _, obs := range durationObservations {
			m.BackupDuration.WithLabelValues("postgres", "full", "local").Observe(obs)
			time.Sleep(15 * time.Second)
		}
	}()

	go func() {
		for i, obs := range sizeObservations {
			m.BackupSize.WithLabelValues("postgres", "full", "local").Observe(obs)
			m.BackupSize.WithLabelValues("mongodb", "full", "local").Observe(x[i])

			time.Sleep(15 * time.Second)
		}
	}()

	restoreDurationObservationsPG := []float64{2.1, 3.4, 5.7, 8.9, 10.2, 1.5, 3, 4.5, 1.5, 7.5}
	restoreDurationObservationsMongo := []float64{1.8, 2.9, 4.1, 7.2, 8.5, 2.3, 3.7, 5.0, 2.1, 6.4}
	restoreDurationObservationsMySQL := []float64{1.2, 2.5, 3.8, 6.1, 7.4, 1.9, 2.7, 4.2, 1.7, 5.8}

	go func() {
		for i, obs := range restoreDurationObservationsPG {
			m.RestoreDuration.WithLabelValues("postgres", "full", "local").Observe(obs)
			m.RestoreDuration.WithLabelValues("mysql", "full", "local").Observe(restoreDurationObservationsMongo[i])
			m.RestoreDuration.WithLabelValues("monog", "full", "local").Observe(restoreDurationObservationsMySQL[i])
			time.Sleep(15 * time.Second)
		}
	}()

	m.RestoreSuccessCount.WithLabelValues("postgres", "full", "local").Add(12)
	m.RestoreFailureCount.WithLabelValues("postgres", "full", "local").Add(7)
	m.RestoreSuccessCount.WithLabelValues("mongodb", "full", "local").Add(2)
	m.RestoreFailureCount.WithLabelValues("mongodb", "full", "local").Add(1)
	m.RestoreSuccessCount.WithLabelValues("mysql", "full", "local").Add(6)
	m.RestoreFailureCount.WithLabelValues("mysql", "full", "local").Add(4)

	y := durationObservations
	slices.Reverse(durationObservations)

	go func() {
		for i, obs := range durationObservations {
			m.ConnectionLatency.WithLabelValues("postgres", "full", "local").Observe(obs)
			m.ConnectionLatency.WithLabelValues("mongodb", "full", "local").Observe(math.Abs(y[i] - obs))
			m.ConnectionLatency.WithLabelValues("mysql", "full", "local").Observe(y[i])

			time.Sleep(15 * time.Second)
		}
	}()

	m.ConnectionFailureCount.WithLabelValues("postgres", "full", "local").Add(7)
	m.ConnectionSuccessCount.WithLabelValues("postgres", "full", "local").Add(12)
	m.ConnectionFailureCount.WithLabelValues("mongodb", "full", "local").Add(5)
	m.ConnectionSuccessCount.WithLabelValues("mongodb", "full", "local").Add(15)
	m.ConnectionFailureCount.WithLabelValues("mysql", "full", "local").Add(9)
	m.ConnectionSuccessCount.WithLabelValues("mysql", "full", "local").Add(11)

}

func NewMetricsWithDummyData(reg prometheus.Registerer) *Metrics {
	m := NewMetrics(reg)
	m.AddDummyData()
	return m
}

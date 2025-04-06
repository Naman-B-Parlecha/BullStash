package backup

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Naman-B-Parlecha/BullStash/internal/config"
	"github.com/jackc/pgx/v5"
)

func PerformBackup(cfg *config.Config) error {
	if err := VerifyPGConfig(cfg); err != nil {
		return fmt.Errorf("postgres configuration error: %w", err)
	}

	connConfig, err := cfg.GetConnConfig()
	if err != nil {
		return fmt.Errorf("failed to get connection config: %w", err)
	}

	if err := takeBaseBackup(connConfig, cfg.Backup.BaseDir); err != nil {
		return fmt.Errorf("base backup failed: %w", err)
	}

	walManager := NewWALManager(connConfig, cfg.Backup.WalDir, cfg.Backup.SlotName)
	go walManager.Start(context.Background())

	return nil
}

func VerifyPGConfig(cfg *config.Config) error {
	connConfig, err := cfg.GetConnConfig()
	if err != nil {
		return err
	}

	conn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	var walLevel, archiveMode string
	if err := conn.QueryRow(context.Background(), "SHOW wal_level").Scan(&walLevel); err != nil {
		return err
	}
	if walLevel != "replica" && walLevel != "logical" {
		return fmt.Errorf("wal_level must be 'replica' or 'logical'")
	}

	if err := conn.QueryRow(context.Background(), "SHOW archive_mode").Scan(&archiveMode); err != nil {
		return err
	}
	if archiveMode != "on" {
		return fmt.Errorf("archive_mode must be 'on'")
	}

	return nil
}

func takeBaseBackup(connConfig *pgx.ConnConfig, backupDir string) error {
	ctx := context.Background()
	conn, err := pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	timestamp := time.Now().Format("20060102-150405")
	backupPath := filepath.Join(backupDir, timestamp)
	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return err
	}

	port := strconv.Itoa(int(connConfig.Port))
	cmd := exec.Command("pg_basebackup",
		"-h", connConfig.Host,
		"-p", port,
		"-U", connConfig.User,
		"-D", backupPath,
		"-Ft",
		"-z",
		"-X", "stream",
		"-P")

	cmd.Env = append(os.Environ(), "PGPASSWORD="+connConfig.Password)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

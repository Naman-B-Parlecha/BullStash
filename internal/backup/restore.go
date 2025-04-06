package backup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Naman-B-Parlecha/BullStash/internal/config"
)

func PerformRestore(cfg *config.Config, targetTime string) error {
	// Parse target time
	_, err := time.Parse(time.RFC3339, targetTime)
	if err != nil {
		return fmt.Errorf("invalid target time format: %w", err)
	}

	// Find most recent base backup
	baseBackup, err := findLatestBaseBackup(cfg.Backup.BaseDir)
	if err != nil {
		return fmt.Errorf("failed to find base backup: %w", err)
	}

	// Stop PostgreSQL if running
	if err := stopPostgreSQL(); err != nil {
		return fmt.Errorf("failed to stop PostgreSQL: %w", err)
	}

	// Clean data directory
	dataDir, err := getPostgreSQLDataDir()
	if err != nil {
		return err
	}
	if err := os.RemoveAll(dataDir); err != nil {
		return fmt.Errorf("failed to clean data directory: %w", err)
	}

	// Restore base backup
	if err := restoreBaseBackup(baseBackup, dataDir); err != nil {
		return fmt.Errorf("failed to restore base backup: %w", err)
	}

	// Configure recovery
	if err := configureRecovery(dataDir, cfg.Backup.WalDir, targetTime); err != nil {
		return fmt.Errorf("failed to configure recovery: %w", err)
	}

	// Start PostgreSQL
	if err := startPostgreSQL(); err != nil {
		return fmt.Errorf("failed to start PostgreSQL: %w", err)
	}

	return nil
}

func findLatestBaseBackup(baseDir string) (string, error) {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return "", err
	}

	var latest time.Time
	var latestPath string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		t, err := time.Parse("20060102-150405", entry.Name())
		if err != nil {
			continue
		}

		if t.After(latest) {
			latest = t
			latestPath = filepath.Join(baseDir, entry.Name())
		}
	}

	if latestPath == "" {
		return "", fmt.Errorf("no valid base backups found")
	}

	return latestPath, nil
}

func restoreBaseBackup(backupPath, dataDir string) error {
	cmd := exec.Command("tar", "-xzf", filepath.Join(backupPath, "base.tar.gz"), "-C", dataDir)
	return cmd.Run()
}

func configureRecovery(dataDir, walDir, targetTime string) error {
	recoveryConf := fmt.Sprintf(`
restore_command = 'cp %s/%%f %%p'
recovery_target_time = '%s'
	`, walDir, targetTime)

	return os.WriteFile(filepath.Join(dataDir, "recovery.conf"), []byte(recoveryConf), 0644)
}

func stopPostgreSQL() error {
	return exec.Command("systemctl", "stop", "postgresql").Run()
}

func startPostgreSQL() error {
	return exec.Command("systemctl", "start", "postgresql").Run()
}

func getPostgreSQLDataDir() (string, error) {
	cmd := exec.Command("psql", "-tAc", "SHOW data_directory")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Postgres struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"postgres"`

	Backup struct {
		BaseDir       string `yaml:"base_dir"`
		WalDir        string `yaml:"wal_dir"`
		SlotName      string `yaml:"slot_name"`
		RetentionDays int    `yaml:"retention_days"`
	} `yaml:"backup"`
}

func LoadConfig(path string) (*Config, error) {
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(home, ".pgbackup.yaml")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) GetConnConfig() (*pgx.ConnConfig, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		c.Postgres.User,
		c.Postgres.Password,
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.Database,
	)

	return pgx.ParseConfig(connStr)
}

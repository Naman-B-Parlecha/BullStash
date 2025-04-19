package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type PostgresConfig struct {
	HOST     string
	PORT     string
	USER     string
	PASSWORD string
	DBNAME   string
}

func GetPostgresConfig() *PostgresConfig {
	godotenv.Load()

	return &PostgresConfig{
		HOST:     GetEnv("POSTGRES_DB_HOST", "localhost"),
		PORT:     GetEnv("POSTGRES_DB_PORT", "5432"),
		USER:     GetEnv("POSTGRES_DB_USER", "postgres"),
		PASSWORD: GetEnv("POSTGRES_DB_PASSWORD", "password"),
		DBNAME:   GetEnv("POSTGRES_DB_NAME", "postgres"),
	}
}

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

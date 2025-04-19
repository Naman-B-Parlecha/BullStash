package config

import "github.com/joho/godotenv"

type MySqlConfig struct {
	HOST     string
	PORT     string
	USER     string
	PASSWORD string
	DBNAME   string
}

func GetMySqlConfig() *MySqlConfig {
	godotenv.Load()

	return &MySqlConfig{
		HOST:     GetEnv("MYSQL_DB_HOST", "localhost"),
		PORT:     GetEnv("MYSQL_DB_PORT", "3306"),
		USER:     GetEnv("MYSQL_DB_USER", "root"),
		PASSWORD: GetEnv("MYSQL_DB_PASSWORD", "password"),
		DBNAME:   GetEnv("MYSQL_DB_NAME", "mysql"),
	}
}

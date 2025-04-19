package config

import (
	"os"
)

func GetEnv(variableName string, defaultValue string) string {
	val, exists := os.LookupEnv(variableName)
	if exists {
		return val
	}
	return defaultValue
}

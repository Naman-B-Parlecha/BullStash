package config

import "github.com/joho/godotenv"

type MongoConfig struct {
	MONGO_URI string
	DBNAME    string
}

func GetMongoConfig() *MongoConfig {
	godotenv.Load()

	return &MongoConfig{
		MONGO_URI: GetEnv("MONGO_URI", ""),
		DBNAME:    GetEnv("MONGO_DB_NAME", ""),
	}
}

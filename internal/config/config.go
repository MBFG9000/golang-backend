package config

import (
	"fmt"
	"os"
	"taskmanager/pkg/modules"

	"github.com/joho/godotenv"
)

func LoadEnviroment() {

	err := godotenv.Load()

	if err != nil {
		panic(err)
	}
}

func GetAuthMiddlewareConfig() *modules.AuthMiddlewareConfig {
	return &modules.AuthMiddlewareConfig{
		ApiKeyHeader: os.Getenv("API_KEY_HEADER"),
		ValidAPIKey:  os.Getenv("VALID_API_KEY"),
	}
}

func GetConfig() *modules.PostgreConfig {

	return &modules.PostgreConfig{
		Host:     os.Getenv("HOST"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("PASSWORD"),
		DBName:   os.Getenv("DATABASE_NAME"),
		SSLMode:  os.Getenv("SSL_MODE"),
	}

}

func GetConnURL(conf *modules.PostgreConfig) string {

	ConnString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.DBName,
		conf.SSLMode,
	)

	return ConnString
}

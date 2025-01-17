package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBAddress              string
	DBName                 string
	JWTExpirationInSeconds int64
	JWTSecret              string
	MAIL_ID                string
	APP_KEY                string
	REDIS_HOST             string
	GEO_LITE_DB            string
}

var EnvConfig = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		PublicHost:             getEnv("PUBLIC_HOST", "http://localhost"),
		Port:                   getEnv("PORT", "3000"),
		DBUser:                 getEnv("DB_USER", "root"),
		DBPassword:             getEnv("DB_PASSWORD", "mypassword"),
		DBAddress:              fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.1"), getEnv("DB_PORT", "3306")),
		DBName:                 getEnv("DB_NAME", "mydb"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXPRATION_IN_SECONDS", 3600*24*7),
		JWTSecret:              getEnv("JWT_SECRET", "i-dont-know-what-to-put-here"),
		MAIL_ID:                getEnv("MAIL_ID", "mymail@mymail.com"),
		APP_KEY:                getEnv("APP_KEY", "appKey"),
		REDIS_HOST:             getEnv("REDIS_HOST", "localhost:6969"),
		GEO_LITE_DB:            getEnv("GEO_LITE_DB", "db/geolite"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		value, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return value
	}
	return fallback
}

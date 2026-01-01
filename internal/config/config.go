package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	MongoURI     string
	MongoDB      string
	JWTSecret    string
	JWTExpiration time.Duration
	CORSOrigin   string
}

func Load() *Config {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using defaults")
	}

	// Parse JWT expiration duration
	jwtExp := getEnv("JWT_EXPIRATION", "24h")
	duration, err := time.ParseDuration(jwtExp)
	if err != nil {
		duration = 24 * time.Hour
	}

	return &Config{
		Port:         getEnv("PORT", "8080"),
		MongoURI:     getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:      getEnv("MONGO_DB", "phone_store"),
		JWTSecret:    getEnv("JWT_SECRET", "default-secret-change-in-production"),
		JWTExpiration: duration,
		CORSOrigin:   getEnv("CORS_ORIGIN", "http://localhost:3000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Application
	AppPort string

	// JWT
	JWTSecret string

	// OpenAI
	OpenAIAPIKey    string
	OpenAIModel     string
	OpenAIMaxTokens int
	OpenAITemp      float32
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "angagrar_db"),

		// Application
		AppPort: getEnv("APP_PORT", "8080"),

		// JWT
		JWTSecret: getEnv("JWT_SECRET", "3RBN1skwbkcF3jp31mVJOuQ0AW38Ut"),

		// OpenAI
		OpenAIAPIKey:    getEnv("OPENAI_API_KEY", ""),
		OpenAIModel:     getEnv("OPENAI_MODEL", "gpt-4o-mini"),
		OpenAIMaxTokens: getEnvInt("OPENAI_MAX_TOKENS", 500),
		OpenAITemp:      getEnvFloat("OPENAI_TEMPERATURE", 0.7),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		if _, err := fmt.Sscanf(value, "%d", &intValue); err == nil {
			return intValue
		}
	}

	return defaultValue
}

func getEnvFloat(key string, defaultValue float32) float32 {
	if value := os.Getenv(key); value != "" {
		var floatValue float32
		if _, err := fmt.Sscanf(value, "%f", &floatValue); err == nil {
			return floatValue
		}
	}

	return defaultValue
}

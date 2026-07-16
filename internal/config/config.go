package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort               string
	DBHost                string
	DBPort                string
	DBUser                string
	DBPassword            string
	DBName                string
	RedisAddr             string
	RedisPassword         string
	RedisDB               string
	RazorpayKeyID         string
	RazorpayKeySecret     string
	RazorpayWebhookSecret string
	JWTSecret             string
	AdminSignupCode       string
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v != "" {
		return v
	}
	return fallback
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found, relying on environment variables")
	}
	return &Config{
		AppPort:               getEnv("APP_PORT", "8080"),
		DBHost:                getEnv("DB_HOST", "localhost"),
		DBPort:                getEnv("DB_PORT", "5432"),
		DBUser:                getEnv("DB_USER", "goshow"),
		DBPassword:            getEnv("DB_PASSWORD", "goshow"),
		DBName:                getEnv("DB_NAME", "goshow"),
		RedisAddr:             getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:         getEnv("REDIS_PASSWORD", ""),
		RedisDB:               getEnv("REDIS_DB", "0"),
		RazorpayKeyID:         getEnv("RAZORPAY_KEY_ID", ""),
		RazorpayKeySecret:     getEnv("RAZORPAY_KEY_SECRET", ""),
		RazorpayWebhookSecret: getEnv("RAZORPAY_WEBHOOK_SECRET", ""),
		JWTSecret:             getEnv("JWTSECRET", ""),
	}

}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}

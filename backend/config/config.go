package config

import "os"

type Config struct {
	Port              string
	RedisURL          string
	FirebaseCredPath  string
	BinanceAPIURL     string
	PriceFetchInterval int
}

func Load() *Config {
	return &Config{
		Port:              getEnv("PORT", "8080"),
		RedisURL:          getEnv("REDIS_URL", "localhost:6379"),
		FirebaseCredPath:  getEnv("FIREBASE_CREDENTIALS", ""),
		BinanceAPIURL:     "https://api.binance.com/api/v3",
		PriceFetchInterval: 10,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}


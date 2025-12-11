package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	RedisURL           string
	BinanceAPIURL      string
	PriceFetchInterval int

	// Telegram Bot 配置
	TelegramBotToken  string // Telegram Bot Token（從 @BotFather 獲得）
	TelegramTestMode  bool   // 測試模式（只 Log 不發送）
	TelegramMyChatID  string // 你自己的 Chat ID（測試用）
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		Port:               getEnv("PORT", "8080"),
		RedisURL:           getEnv("REDIS_URL", "localhost:6379"),
		BinanceAPIURL:      "https://api.binance.com/api/v3",
		PriceFetchInterval: 10,

		// Telegram 配置
		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		TelegramTestMode: getEnvBool("TELEGRAM_TEST_MODE", true), // 預設測試模式
		TelegramMyChatID: getEnv("TELEGRAM_MY_CHAT_ID", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		b, err := strconv.ParseBool(value)
		if err == nil {
			return b
		}
	}
	return defaultValue
}

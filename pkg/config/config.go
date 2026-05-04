package config

import "os"

type Config struct {
	Port         string
	GeminiAPIKey string
	XenditAPIKey string
}

func New() *Config {
	return &Config{
		Port:         os.Getenv("PORT"),
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
		XenditAPIKey: os.Getenv("XENDIT_API_KEY"),
	}
}

package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port                string
	BufUploadSizeInfo   int64
	BufUploadSizeCreate int64
	BufUploadSizeMail   int64
	SMTPPort            string
	SMTPHost            string
	SMTPUser            string
	SMTPPassword        string
}

// New returns a new Config struct
func New() *Config {
	return &Config{
		Port:                getEnv("PORT", ":3000"),
		BufUploadSizeInfo:   int64(getEnvAsInt("MAX_UPLOAD_SIZE_INFO", 10485760)),
		BufUploadSizeCreate: int64(getEnvAsInt("MAX_UPLOAD_SIZE_CREATE", 33554432)),
		BufUploadSizeMail:   int64(getEnvAsInt("MAX_UPLOAD_SIZE_MAIL", 10485760)),
		SMTPHost:            getEnv("SMTP_HOST", ""),
		SMTPPort:            getEnv("SMTP_PORT", "587"),
		SMTPUser:            getEnv("SMTP_USER", ""),
		SMTPPassword:        getEnv("SMTP_PASSWORD", ""),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// Helper to read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

// Helper to read an environment variable into a string slice or return default value
func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}
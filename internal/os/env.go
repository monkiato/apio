package os

import (
	"os"
	"strconv"
)

//GetEnv look for environment variable
func GetEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}

// GetIntEnv look for environment variable and convert to int
func GetIntEnv(key string, defaultValue int) int {
	if value, ok := os.LookupEnv(key); ok {
		if val, err := strconv.Atoi(value); err == nil {
			return val
		}
	}

	return defaultValue
}

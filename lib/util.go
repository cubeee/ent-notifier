package lib

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetRequiredEnv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		panic(fmt.Errorf("required environment variable missing: %s", key))
	}
	return value
}

func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func GetEnvInt(key string, fallback int) int {
	value := GetEnv(key, strconv.Itoa(fallback))
	intValue, err := strconv.Atoi(value)
	if err != nil {
		fmt.Println("Failed to read env var into an int")
		panic(err)
	}
	return intValue
}

func GetEnvList(key string, delimiter string) []string {
	value := GetEnv(key, "")
	return strings.Split(value, delimiter)
}

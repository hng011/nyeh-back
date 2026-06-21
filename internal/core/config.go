package core

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	// APP
	ENV   string
	PORT  uint
	DEBUG bool

	// AUTH
	JWT_AUTH_TOKEN           string
	TOKEN_DIGEST             string
	TTL_ACCESS_TOKEN_MINUTES uint
	TTL_REFRESH_TOKEN_HOURS  uint

	// GOOGLE SSO
	GOOGLE_ALLOWED_EMAIL       string
	GOOGLE_OAUTH_CLIENT_ID     string
	GOOGLE_OAUTH_CLIENT_SECRET string

	// REDIS
	REDIS_ADDR     string
	REDIS_PASSWORD string
}

var Settings *Config // Global Settings

func LoadEnv() {
	file, err := os.Open(".env")
	if err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			if len(line) == 0 || strings.HasPrefix(line, "#") {
				continue
			}

			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key, val := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
				val = strings.Trim(val, `"'`)
				os.Setenv(key, val)
			}
		}
		if err := scanner.Err(); err != nil {
			log.Panicf("[Panic] error reading .env file: %v", err)
		}
	} else {
		log.Println("No .env file found.")
	}

	Settings = &Config{
		ENV:                        getEnv[string]("ENV"),
		PORT:                       getEnv[uint]("PORT", 8080),
		DEBUG:                      getEnv[bool]("DEBUG", false),
		JWT_AUTH_TOKEN:             getEnv[string]("JWT_AUTH_TOKEN"),
		TOKEN_DIGEST:               getEnv[string]("TOKEN_DIGEST"),
		GOOGLE_ALLOWED_EMAIL:       getEnv[string]("GOOGLE_ALLOWED_EMAIL"),
		GOOGLE_OAUTH_CLIENT_ID:     getEnv[string]("GOOGLE_OAUTH_CLIENT_ID"),
		GOOGLE_OAUTH_CLIENT_SECRET: getEnv[string]("GOOGLE_OAUTH_CLIENT_SECRET"),
		REDIS_ADDR:                 getEnv[string]("REDIS_ADDR"),
		REDIS_PASSWORD:             getEnv[string]("REDIS_PASSWORD"),
		TTL_ACCESS_TOKEN_MINUTES:   getEnv[uint]("TTL_ACCESS_TOKEN_MINUTES"),
		TTL_REFRESH_TOKEN_HOURS:    getEnv[uint]("TTL_REFRESH_TOKEN_HOURS"),
	}
}

func getEnv[T uint | bool | string](key string, fallback ...T) T {
	value, exists := os.LookupEnv(key)

	if !exists {
		if len(fallback) > 0 {
			return fallback[0]
		}
		var def T
		return def
	}

	var targetType any = *new(T)
	switch targetType.(type) {
	case string:
		var res any = value
		return res.(T)

	case uint:
		parsed, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			if len(fallback) > 0 {
				return fallback[0]
			}
			var def T
			return def
		}
		var res any = uint(parsed)
		return res.(T)

	case bool:
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			if len(fallback) > 0 {
				return fallback[0]
			}
			var zero T
			return zero
		}
		var res any = parsed
		return res.(T)
	}

	var def T
	return def
}

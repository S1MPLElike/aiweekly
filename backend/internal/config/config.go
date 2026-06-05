package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string

	MySQLDSN string

	JWTSecret string

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration

	CookieDomain string
	CookieSecure bool

	CORSAllowOrigins []string

	DeepseekAPIKey string
}

func Load() Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return Config{
		Port:             getEnv("APP_PORT", "8082"),
		MySQLDSN:         mustEnv("MYSQL_DSN"),
		JWTSecret:        mustEnv("JWT_SECRET"),
		AccessTokenTTL:   getEnvDuration("ACCESS_TOKEN_TTL", 2*time.Hour),
		RefreshTokenTTL:  getEnvDuration("REFRESH_TOKEN_TTL", 7*24*time.Hour),
		CookieDomain:     getEnv("COOKIE_DOMAIN", ""),
		CookieSecure:     getEnvBool("COOKIE_SECURE", false),
		CORSAllowOrigins: getEnvCSV("CORS_ALLOW_ORIGINS", "http://localhost:5173"),
		DeepseekAPIKey:   getEnv("DEEPSEEK_API_KEY", ""),
	}
}

func mustEnv(key string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		panic("missing required env: " + key)
	}
	return v
}

func getEnv(key, def string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	return v
}

func getEnvBool(key string, def bool) bool {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}

func getEnvDuration(key string, def time.Duration) time.Duration {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}
	return d
}

func getEnvCSV(key string, def string) []string {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		raw = def
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

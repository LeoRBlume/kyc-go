package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	// App
	AppEnv   string
	HTTPPort string

	// DB
	DBDriver string
	DBPath   string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	DBTimeZone string
}

func Load() (*Config, error) {
	cfg := &Config{
		// App
		AppEnv:   getenv("APP_ENV", "local"),
		HTTPPort: getenv("HTTP_PORT", getenv("PORT", "8080")),

		// DB
		DBDriver: getenv("DB_DRIVER", ""),
		DBPath:   getenv("DB_PATH", ""),

		DBHost:     getenv("DB_HOST", ""),
		DBPort:     getenv("DB_PORT", ""),
		DBUser:     getenv("DB_USER", ""),
		DBPassword: getenv("DB_PASSWORD", ""),
		DBName:     getenv("DB_NAME", ""),
		DBSSLMode:  getenv("DB_SSLMODE", "disable"),
		DBTimeZone: getenv("DB_TIMEZONE", "UTC"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if strings.TrimSpace(c.HTTPPort) == "" {
		return fmt.Errorf("HTTP_PORT is required (or PORT)")
	}

	switch strings.ToLower(c.DBDriver) {
	case "sqlite":
		if c.DBPath == "" {
			return fmt.Errorf("DB_PATH is required when DB_DRIVER=sqlite")
		}
		return nil

	case "postgres":
		missing := []string{}
		if c.DBHost == "" {
			missing = append(missing, "DB_HOST")
		}
		if c.DBPort == "" {
			missing = append(missing, "DB_PORT")
		}
		if c.DBUser == "" {
			missing = append(missing, "DB_USER")
		}
		if c.DBPassword == "" {
			missing = append(missing, "DB_PASSWORD")
		}
		if c.DBName == "" {
			missing = append(missing, "DB_NAME")
		}
		if len(missing) > 0 {
			return fmt.Errorf("missing env vars for postgres: %v", missing)
		}
		return nil

	default:
		return fmt.Errorf("unsupported DB_DRIVER: %s (expected sqlite or postgres)", c.DBDriver)
	}
}

func getenv(key, def string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	return v
}

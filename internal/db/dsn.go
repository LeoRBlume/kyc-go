package db

import (
	"fmt"
	"strings"

	"kyc-sim/internal/config"
)

func BuildDSN(cfg *config.Config) (string, error) {
	switch strings.ToLower(cfg.DBDriver) {
	case "sqlite":
		return cfg.DBPath, nil

	case "postgres":
		return fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBName,
			cfg.DBSSLMode,
			cfg.DBTimeZone,
		), nil

	default:
		return "", fmt.Errorf("unsupported DB_DRIVER: %s", cfg.DBDriver)
	}
}

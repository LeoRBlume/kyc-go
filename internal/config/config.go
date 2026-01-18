package config

type Config struct {
	AppEnv     string
	HTTPPort   string
	DBDriver   string
	DBDsn      string
	DBLogLevel string
}

func Load() Config {
	// Defaults adequados para dev local
	return Config{
		AppEnv:     getEnv("APP_ENV", "local"),
		HTTPPort:   getEnv("HTTP_PORT", "8080"),
		DBDriver:   getEnv("DB_DRIVER", "sqlite"),
		DBDsn:      getEnv("DB_DSN", "./data/kyc.db"),
		DBLogLevel: getEnv("DB_LOG_LEVEL", "warn"),
	}
}

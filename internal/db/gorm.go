package db

import (
	"errors"
	"kyc-sim/internal/config"
	"log"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

func NewGormDB(cfg config.Config) (*gorm.DB, error) {
	if cfg.DBDriver != "sqlite" {
		return nil, errors.New("only sqlite is supported in this stage")
	}

	// Log level
	var level glogger.LogLevel
	switch cfg.DBLogLevel {
	case "silent":
		level = glogger.Silent
	case "error":
		level = glogger.Error
	case "info":
		level = glogger.Info
	default:
		level = glogger.Warn
	}

	logger := glogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		glogger.Config{
			SlowThreshold: time.Second,
			LogLevel:      level,
			Colorful:      true,
		},
	)

	db, err := gorm.Open(sqlite.Open(cfg.DBDsn), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		return nil, err
	}

	// PRAGMAs úteis para concorrência (WAL) e tempo de espera
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	_, _ = sqlDB.Exec("PRAGMA journal_mode=WAL;")
	_, _ = sqlDB.Exec("PRAGMA synchronous=NORMAL;")
	_, _ = sqlDB.Exec("PRAGMA busy_timeout=5000;")
	_, _ = sqlDB.Exec("PRAGMA foreign_keys=ON;")

	return db, nil
}

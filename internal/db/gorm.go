package db

import (
	"fmt"
	"kyc-sim/internal/config"
	"log"
	"os"
	"strings"
	"time"

	"github.com/glebarez/sqlite" // mantém: sqlite sem CGO
	"gorm.io/driver/postgres"    // adiciona: postgres
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

func NewGormDB(cfg config.Config) (*gorm.DB, error) {
	logger := glogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		glogger.Config{
			SlowThreshold: time.Second,
			Colorful:      true,
		},
	)

	dsn, err := BuildDSN(&cfg)
	if err != nil {
		return nil, err
	}

	var db *gorm.DB

	switch strings.ToLower(cfg.DBDriver) {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{Logger: logger})
		if err != nil {
			return nil, err
		}

		sqlDB, e := db.DB()
		if e != nil {
			return nil, e
		}
		_, _ = sqlDB.Exec("PRAGMA journal_mode=WAL;")
		_, _ = sqlDB.Exec("PRAGMA synchronous=NORMAL;")
		_, _ = sqlDB.Exec("PRAGMA busy_timeout=5000;")
		_, _ = sqlDB.Exec("PRAGMA foreign_keys=ON;")

	case "postgres":
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger})
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported DB_DRIVER: %s (expected sqlite|postgres)", cfg.DBDriver)
	}

	return db, nil
}

package db

import (
	"task-5/internal/config"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(cfg config.DBConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	gormConfig := &gorm.Config{
		TranslateError: true,
		Logger:         logger.Default.LogMode(logger.Silent),
	}

	switch cfg.Driver {
	case "postgres":
		db, err = gorm.Open(postgres.Open(cfg.DSN), gormConfig)
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(cfg.DSN), gormConfig)
	default:
		return nil, &DatabaseError{DSN: cfg.DSN, Err: gorm.ErrUnsupportedDriver}
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}

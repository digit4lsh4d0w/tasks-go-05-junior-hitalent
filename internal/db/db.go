package db

import (
	"task-5/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDatabase(cfg config.DBConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	switch cfg.Driver {
	case "postgres":
		db, err = gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{TranslateError: true})
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(cfg.DSN), &gorm.Config{TranslateError: true})
	default:
		return nil, &DatabaseError{DSN: cfg.DSN, Err: gorm.ErrUnsupportedDriver}
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}

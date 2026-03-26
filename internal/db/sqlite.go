package db

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"xhs-go-cli/internal/model"
)

type DB struct {
	*gorm.DB
}

func Open(path string) (*DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := migrate(db); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (d *DB) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("get underlying sql.DB: %w", err)
	}
	return sqlDB.Close()
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Source{},
		&model.GeneratedQuery{},
		&model.SearchResult{},
		&model.Detail{},
		&model.Qualification{},
	)
}

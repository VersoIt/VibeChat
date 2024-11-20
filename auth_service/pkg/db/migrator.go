package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

type Migrator struct {
	db             *sqlx.DB
	migrationsPath string
	dialect        string
}

func NewMigrator(db *sqlx.DB, migrationsPath, dialect string) *Migrator {
	return &Migrator{db: db, migrationsPath: migrationsPath, dialect: dialect}
}

func (m *Migrator) Migrate() error {
	if err := goose.SetDialect(m.dialect); err != nil {
		return err
	}

	if err := goose.Up(m.db.DB, m.migrationsPath); err != nil {
		return err
	}

	return nil
}

package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Name     string
	Host     string
	Port     string
	User     string
	Password string
	SSLMode  string
}

func NewPostgresDb(cfg Config) (*sqlx.DB, error) {
	dsn := getDsn(cfg)
	logrus.Info("db dsn: ", dsn)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func getDsn(cfg Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)
}

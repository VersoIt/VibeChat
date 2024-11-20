package main

import (
	"Messenger/internal/di"
	"Messenger/internal/service"
	"Messenger/pkg/db"
	http2 "Messenger/pkg/http"
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		logrus.Error("error loading .env file")
		return
	}

	database, err := db.NewPostgresDb(db.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		Name:     os.Getenv("DB_NAME"),
	})
	if err != nil {
		logrus.Error(err)
	}
	defer func(database *sqlx.DB) {
		err = database.Close()
		if err != nil {
			logrus.Error(err)
		}
	}(database)

	migrator := db.NewMigrator(database, "./schema", "postgres")
	if err = migrator.Migrate(); err != nil {
		logrus.Error(err)
		return
	}

	refreshTokenTTL, err := parseDurationFromEnv("REFRESH_TTL")
	if err != nil {
		logrus.Error(err)
		return
	}

	accessTokenTTL, err := parseDurationFromEnv("ACCESS_TTL")
	if err != nil {
		logrus.Error(err)
		return
	}

	h, err := di.InitializeAuthHandler(database, service.AuthCfg{
		AccessTokenSecret:  []byte(os.Getenv("ACCESS_TOKEN_SECRET")),
		RefreshTokenSecret: []byte(os.Getenv("ACCESS_TOKEN_SECRET")),
		PasswordSecret:     []byte(os.Getenv("PASSWORD_SECRET")),
		RefreshTokenTTL:    refreshTokenTTL,
		AccessTokenTTL:     accessTokenTTL,
	})
	if err != nil {
		logrus.Error(err)
		return
	}

	server := http2.NewServer(":8080", h)
	go func() {
		if err = server.Run(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				logrus.Error(err)
			} else {
				logrus.Info("server closed")
			}
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-exit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = server.Shutdown(ctx); err != nil {
		logrus.Error(err)
		return
	}

	logrus.Info("app closed gracefully")
}

func parseDurationFromEnv(key string) (time.Duration, error) {
	valueStr := os.Getenv(key)
	duration, err := time.ParseDuration(valueStr)
	if err != nil {
		return 0, err
	}
	return duration, nil
}

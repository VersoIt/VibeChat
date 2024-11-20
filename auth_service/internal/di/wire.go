//go:build wireinject
// +build wireinject

package di

import (
	"Messenger/internal/handler/http"
	"Messenger/internal/repository"
	"Messenger/internal/service"
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
)

func initializeAuthHandler(db *sqlx.DB, authCfg service.AuthCfg) (*http.Handler, error) {
	wire.Build(
		repository.NewUserRepository,
		repository.NewBlockedRefreshTokenRepository,
		wire.Bind(new(service.UserRepository), new(*repository.UserRepository)),
		wire.Bind(new(service.BlockedRefreshTokenRepository), new(*repository.BlockedRefreshTokenRepository)),
		service.NewAuth,
		http.NewHandler,
	)
	return new(http.Handler), nil
}

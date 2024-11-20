package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type BlockedRefreshTokenRepository struct {
	db *sqlx.DB
}

func NewBlockedRefreshTokenRepository(db *sqlx.DB) *BlockedRefreshTokenRepository {
	return &BlockedRefreshTokenRepository{db: db}
}

func (r *BlockedRefreshTokenRepository) Create(token string) error {
	query := fmt.Sprintf("INSERT INTO %s(token) VALUES($1)", blockedRefreshTokensTableName)
	if _, err := r.db.Exec(query, token); err != nil {
		return err
	}
	return nil
}

func (r *BlockedRefreshTokenRepository) Exists(token string) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE token=$1", blockedRefreshTokensTableName)
	var count int
	err := r.db.Get(&count, query, token)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *BlockedRefreshTokenRepository) Delete(token string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE token=$1", blockedRefreshTokensTableName)
	if _, err := r.db.Exec(query, token); err != nil {
		return err
	}
	return nil
}

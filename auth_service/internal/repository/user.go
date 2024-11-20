package repository

import (
	"Messenger/internal/model"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

const (
	usersTableName                = "users"
	blockedRefreshTokensTableName = "blocked_refresh_tokens"
)

func (r *UserRepository) CreateUser(user model.User) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s(login, email, password) VALUES($1, $2, $3)", usersTableName)
	if err := r.db.QueryRow(query, user.Login, user.Email, user.Password).Scan(&user.Id); err != nil {
		return 0, err
	}
	return user.Id, nil
}

func (r *UserRepository) GetUserByEmail(email, password string) (model.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE email = $1 AND password = $2", usersTableName)
	var user model.User
	if err := r.db.Get(&user, query, email, password); err != nil {
		return user, err
	}
	return user, nil
}

func (r *UserRepository) GetUserById(id int) (model.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", usersTableName)
	var user model.User
	if err := r.db.Get(&user, query, id); err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByLogin(login, password string) (model.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE login = $1 AND password = $2", usersTableName)
	var user model.User
	if err := r.db.Get(&user, query, login, password); err != nil {
		return user, err
	}
	return user, nil
}

func (r *UserRepository) GetAllUsers() ([]model.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s", usersTableName)
	var users []model.User
	if err := r.db.Select(&users, query); err != nil {
		return nil, err
	}
	return users, nil
}

package service

import (
	"Messenger/internal/model"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type AuthCfg struct {
	AccessTokenSecret, RefreshTokenSecret, PasswordSecret []byte
	AccessTokenTTL, RefreshTokenTTL                       time.Duration
}

type UserRepository interface {
	CreateUser(user model.User) (int, error)
	GetUserByEmail(email, password string) (model.User, error)
	GetUserById(id int) (model.User, error)
	GetUserByLogin(login, password string) (model.User, error)
	GetAllUsers() ([]model.User, error)
}

type BlockedRefreshTokenRepository interface {
	Create(token string) error
	Exists(token string) (bool, error)
	Delete(token string) error
}

var (
	ErrTokenExpired = errors.New("token expired")
	ErrTokenInvalid = errors.New("token invalid")
)

type Auth struct {
	userRepo                UserRepository
	blockedRefreshTokenRepo BlockedRefreshTokenRepository
	accessSecret            []byte
	refreshSecret           []byte
	passwordSecret          []byte
	accessTokenTTL          time.Duration
	refreshTokenTTL         time.Duration
}

func NewAuth(userRepo UserRepository, blockedRefreshTokenRepo BlockedRefreshTokenRepository, authCfg AuthCfg) *Auth {
	return &Auth{userRepo: userRepo, blockedRefreshTokenRepo: blockedRefreshTokenRepo, accessSecret: authCfg.AccessTokenSecret, refreshSecret: authCfg.RefreshTokenSecret, passwordSecret: authCfg.PasswordSecret, accessTokenTTL: authCfg.AccessTokenTTL, refreshTokenTTL: authCfg.RefreshTokenTTL}
}

type claims struct {
	userId int
	jwt.RegisteredClaims
}

func (a *Auth) CreateUser(user model.User) (int, error) {
	user.Password = a.generatePasswordHash(user.Password)
	return a.userRepo.CreateUser(user)
}

func (a *Auth) GetUserById(id int) (model.User, error) {
	return a.userRepo.GetUserById(id)
}

func (a *Auth) GetUserByLogin(login, password string) (model.User, error) {
	return a.userRepo.GetUserByLogin(login, password)
}

func (a *Auth) GetUserByEmail(email, password string) (model.User, error) {
	return a.userRepo.GetUserByEmail(email, a.generatePasswordHash(password))
}

func (a *Auth) GenerateAuthTokensByEmail(email, password string) (model.AuthTokens, error) {
	user, err := a.GetUserByEmail(email, password)
	if err != nil {
		return model.AuthTokens{}, err
	}

	return a.generateAuthTokensByUser(user)
}

func (a *Auth) generateAuthTokensByUser(user model.User) (model.AuthTokens, error) {

	refreshToken, err := a.GenerateRefreshToken(user.Id)
	if err != nil {
		return model.AuthTokens{}, err
	}
	accessToken, err := a.GenerateAccessToken(user.Id)
	if err != nil {
		return model.AuthTokens{}, err
	}

	return model.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *Auth) GenerateAuthTokensByLogin(login, password string) (model.AuthTokens, error) {
	user, err := a.GetUserByLogin(login, a.generatePasswordHash(password))
	if err != nil {
		return model.AuthTokens{}, err
	}

	return a.generateAuthTokensByUser(user)
}

func (a *Auth) GenerateAccessToken(id int) (string, error) {
	return a.generateToken(id, a.accessTokenTTL, a.accessSecret)
}

func (a *Auth) GenerateRefreshToken(id int) (string, error) {
	return a.generateToken(id, a.refreshTokenTTL, a.refreshSecret)
}

func (a *Auth) generateToken(id int, ttl time.Duration, secret []byte) (string, error) {
	c := &claims{
		userId: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(secret)
}

func (a *Auth) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum(a.passwordSecret))
}

func (a *Auth) SignOut(refreshToken string) error {
	return a.blockedRefreshTokenRepo.Create(refreshToken)
}

func (a *Auth) ValidateAccessToken(token string) (bool, error) {
	_, err := parseJwt(token, a.accessSecret)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (a *Auth) ValidateRefreshToken(token string) (bool, error) {
	blocked, err := a.blockedRefreshTokenRepo.Exists(token)
	if err != nil {
		return false, err
	}
	if blocked {
		return false, nil
	}

	_, err = parseJwt(token, a.refreshSecret)
	if err != nil {

		if errors.Is(err, ErrTokenExpired) {
			_ = a.blockedRefreshTokenRepo.Delete(token)
		}

		return false, err
	}

	return true, nil
}

func (a *Auth) GenerateAccessTokenFromRefreshToken(refresh string) (string, error) {
	c, err := parseJwt(refresh, a.refreshSecret)
	if err != nil {
		return "", err
	}
	return a.GenerateAccessToken(c.userId)
}

func parseJwt(token string, secret []byte) (*claims, error) {
	t, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}
	c, ok := t.Claims.(*claims)
	if !ok || !t.Valid {
		return nil, ErrTokenInvalid
	}

	if c.ExpiresAt != nil && time.Now().After(c.ExpiresAt.Time) {
		return nil, ErrTokenExpired
	}

	return c, nil
}

func (a *Auth) GetAllUsers() ([]model.User, error) {
	return a.userRepo.GetAllUsers()
}

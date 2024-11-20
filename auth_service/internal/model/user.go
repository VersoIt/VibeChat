package model

type User struct {
	Id       int    `json:"id" db:"id"`
	Login    string `json:"login" db:"login"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

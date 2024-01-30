package models

type Account interface {
	PasswordMatching(password string) bool
	UpdatePassword(password string)
	CreateSession() string
}

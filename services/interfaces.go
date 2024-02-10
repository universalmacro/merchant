package services

type Account interface {
	MerchantId() uint
	PasswordMatching(password string) bool
	UpdatePassword(password string)
	CreateSession() string
}

package services

type Account interface {
	ID() uint
	MerchantId() uint
	PasswordMatching(password string) bool
	UpdatePassword(password string)
	CreateSession() string
}

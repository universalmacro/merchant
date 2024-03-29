package services

import (
	"github.com/golang-jwt/jwt"
	"github.com/universalmacro/common/snowflake"
)

var sessionIdGenerator = snowflake.NewIdGenertor(0)

type Claims struct {
	jwt.StandardClaims
	ID         string `json:"id"`
	MerchantID string `json:"merchantId"`
	Account    string `json:"account"`
	Type       string `json:"type"`
}

type MemberClaims struct {
	jwt.StandardClaims
	ID         string `json:"id"`
	MerchantID string `json:"merchantId"`
	Type       string `json:"type"`
}

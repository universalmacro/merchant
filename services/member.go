package services

import (
	"github.com/universalmacro/common/utils"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/ioc"
)

type Member struct {
	entities.Member
}

func (m *Member) GenerateToken() string {
	jwtSigner := ioc.GetJwtSigner()
	token, _ := jwtSigner.SignJwt(MemberClaims{
		ID:         utils.UintToString(m.ID),
		MerchantID: utils.UintToString(m.MerchantId),
		Type:       "MEMBER",
	})
	return token
}

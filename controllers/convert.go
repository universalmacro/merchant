package controllers

import (
	api "github.com/universalmacro/merchant-api-interfaces"
	"github.com/universalmacro/merchant/services/models"
)

func ConvertMerchant(m models.Merchant) api.Merchant {
	return api.Merchant{
		Id:              m.StringID(),
		ShortMerchantId: m.ShortMerchantId(),
		Account:         m.Account(),
		UpdatedAt:       m.UpdatedAt().Unix(),
		CreatedAt:       m.CreatedAt().Unix(),
		Name:            m.Name(),
	}
}

func ConvertSpace(m models.Space) api.Space {
	return api.Space{
		Id:   m.StringID(),
		Name: m.Name,
	}
}

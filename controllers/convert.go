package controllers

import (
	"github.com/universalmacro/common/utils"
	api "github.com/universalmacro/merchant-api-interfaces"
	"github.com/universalmacro/merchant/services"
)

func ConvertMerchant(m services.Merchant) api.Merchant {
	return api.Merchant{
		Id:              m.StringID(),
		ShortMerchantId: m.ShortMerchantId(),
		Account:         m.Account(),
		UpdatedAt:       m.UpdatedAt().Unix(),
		CreatedAt:       m.CreatedAt().Unix(),
		Name:            m.Name(),
	}
}

func ConvertSpace(m *services.Space) api.Space {
	return api.Space{
		Id:   m.StringID(),
		Name: m.Name(),
	}
}

func ConvertTable(t *services.Table) api.Table {
	return api.Table{
		Id:    utils.UintToString(t.ID()),
		Label: t.Label(),
	}
}

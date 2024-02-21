package controllers

import (
	"github.com/universalmacro/common/convert"
	"github.com/universalmacro/common/utils"
	api "github.com/universalmacro/merchant-api-interfaces"
	"github.com/universalmacro/merchant/dao/entities"
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

func ConvertFood(f *services.Food) api.Food {
	return api.Food{
		Id:          f.StringID(),
		Name:        f.Name(),
		Categories:  f.Categories(),
		Description: f.Description(),
		Price:       f.Price(),
		Attributes:  ConvertFoodAttributes(f.Attributes()),
		FixedOffset: f.FixedOffset(),
		Status:      api.FoodStatus(f.Status()),
		UpdatedAt:   f.UpdatedAt.Unix(),
		CreatedAt:   f.CreatedAt.Unix(),
	}
}

func ConvertFoodAttributes(f entities.Attributes) []api.FoodAttribute {
	attributes := make([]api.FoodAttribute, len(f))
	for i := range f {
		var options []api.FoodAttributesOption
		convert.ConvertByJSON(f[i].Options, &options)
		attributes[i] = api.FoodAttribute{
			Label:   f[i].Label,
			Options: options,
		}
	}
	return attributes
}

func ConvertPrinter(p *services.Printer) api.Printer {
	return api.Printer{
		Id:    p.StringID(),
		Name:  p.Name,
		Sn:    p.Sn,
		Type:  api.PrinterType(p.Type),
		Model: api.PrinterModel(p.Printer.Model),
	}
}

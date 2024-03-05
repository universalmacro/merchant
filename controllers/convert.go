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
		Image:       f.Image(),
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

func ConvertFoodSpec(f *services.FoodSpec) api.FoodSpec {
	var spec []api.Spec
	if f.Spec.Len() != 0 {
		for _, s := range f.Spec.Spec {
			spec = append(spec, api.Spec{
				Attribute: s.Attribute,
				Optioned:  s.Optioned,
			})
		}
	}
	return api.FoodSpec{
		Food: ConvertFood(f.Food),
		Spec: &spec,
	}
}

func ConvertOrder(o *services.Order) api.Order {
	if o == nil {
		return api.Order{}
	}
	// var foods []api.Food
	// for i := range o.Foods {
	// 	foods = append(foods, ConvertFood(&o.Foods[i].Food))
	// }
	return api.Order{
		Id:         o.StringID(),
		UpdatedAt:  o.UpdatedAt.Unix(),
		CreatedAt:  o.CreatedAt.Unix(),
		Status:     api.OrderStatus(o.Status),
		TableLabel: o.TableLabel,
		Code:       o.Code(),
	}
}

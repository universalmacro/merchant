package factories

import (
	"github.com/universalmacro/common/utils"
	api "github.com/universalmacro/merchant-api-interfaces"
	"github.com/universalmacro/merchant/services"
)

func NewFoodSpec(foodSpec api.FoodSpec) *services.FoodSpec {
	foodService := services.GetFoodService()
	food := foodService.GetById(utils.StringToUint(foodSpec.Food.Id))
	if food == nil {
		return nil
	}
	specMap := make(map[string]string)
	if foodSpec.Spec == nil {
		return &services.FoodSpec{Food: food}
	}
	for _, spec := range *foodSpec.Spec {
		specMap[spec.Attribute] = spec.Optioned
	}
	f := food.CreateFoodSpec(specMap)
	return &f
}

func NewFoodSpecs(foodSpecs []api.FoodSpec) []services.FoodSpec {
	var result []services.FoodSpec
	for _, foodSpec := range foodSpecs {
		if food := NewFoodSpec(foodSpec); food != nil {
			result = append(result, *food)
		}
	}
	return result
}

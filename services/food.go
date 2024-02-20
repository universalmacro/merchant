package services

import (
	"errors"

	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/common/utils"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/dao/repositories"
)

func NewFood() *Food {
	return &Food{
		Food: &entities.Food{},
	}
}

type Food struct {
	*entities.Food
}

func (f *Food) ID() uint {
	return f.Food.ID
}

func (f *Food) StringID() string {
	return utils.UintToString(f.ID())
}

func (f *Food) SetName(name string) *Food {
	f.Food.Name = name
	return f
}

func (f *Food) Name() string {
	return f.Food.Name
}

func (f *Food) Status() string {
	return f.Food.Status
}

func (f *Food) SetStatus(status string) *Food {
	f.Food.Status = status
	return f
}

func (f *Food) Categories() []string {
	return f.Food.Categories
}

func (f *Food) SetCategories(categories []string) *Food {
	f.Food.Categories = categories
	return f
}

func (f *Food) SetDescription(description string) *Food {
	f.Food.Description = description
	return f
}

func (f *Food) Description() string {
	return f.Food.Description
}

func (f *Food) SetPrice(price int64) *Food {
	f.Food.Price = price
	return f
}

func (f *Food) Price() int64 {
	return f.Food.Price
}

func (f *Food) SetFixedOffset(fixedOffset *int64) *Food {
	f.Food.FixedOffset = fixedOffset
	return f
}

func (f *Food) FixedOffset() *int64 {
	return f.Food.FixedOffset
}

func (f *Food) SetImage(image string) *Food {
	f.Food.Image = image
	return f
}

func (f *Food) Image() string {
	return f.Food.Image
}

func (f *Food) SetAttributes(attributes entities.Attributes) *Food {
	f.Food.Attributes = attributes
	return f
}

func (f *Food) Attributes() entities.Attributes {
	return f.Food.Attributes
}

func (f *Food) AddAttribute(label string, options ...entities.Option) (*Food, error) {
	for _, attr := range f.Food.Attributes {
		if attr.Label == label {
			return nil, errors.New("attribute label duplicated")
		}
	}
	optionsMap := make(map[string]struct{})
	for _, option := range options {
		if _, ok := optionsMap[option.Label]; ok {
			return nil, errors.New(label + " attribute option " + option.Label + " label duplicated")
		}
		optionsMap[option.Label] = struct{}{}
	}
	f.Food.Attributes = append(f.Food.Attributes, entities.Attribute{
		Label:   label,
		Options: options,
	})
	return f, nil
}

func (f *Food) Submit() *Food {
	repo := repositories.GetFoodRepository()
	repo.Save(f.Food)
	return f
}

func (r *Food) Create() *Food {
	repo := repositories.GetFoodRepository()
	repo.Create(r.Food)
	return r
}

func (f *Food) Delete() {
	repo := repositories.GetFoodRepository()
	repo.Delete(f.Food)
}

func (f *Food) Space() *Space {
	return newSpaceService().GetSpace(f.SpaceID)
}

func (f *Food) Granted(account Account) bool {
	return f.Space().Granted(account)
}

var foodServiceSingleton = singleton.SingletonFactory(newFoodService, singleton.Eager)

func GetFoodService() *FoodService {
	return foodServiceSingleton.Get()
}

func newFoodService() *FoodService {
	return &FoodService{
		foodRepository: repositories.GetFoodRepository(),
	}
}

type FoodService struct {
	foodRepository *repositories.FoodRepository
}

func (self *FoodService) GetById(id uint) *Food {
	f, _ := self.foodRepository.GetById(id)
	if f == nil {
		return nil
	}
	return &Food{Food: f}
}

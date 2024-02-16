package services

import (
	"github.com/universalmacro/common/utils"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/dao/repositories"
)

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

func (f *Food) AddAttribute(label string, options ...Option) *Food {
	f.Food.Attributes = append(f.Food.Attributes, entities.Attribute{
		Label: label,
	})
	return f
}

func (f *Food) Submit() *Food {
	repo := repositories.GetFoodRepository()
	repo.Save(f.Food)
	return f
}

func (f *Food) Delete() {
	repo := repositories.GetFoodRepository()
	repo.Delete(f.Food)
}

func (f *Food) Space() *Space {
	return newSpaceService().GetSpace(f.SpaceID)
}

type Option struct {
	Label   string   `json:"label"`
	Options []Option `json:"options"`
}

func (o *Option) SetLabel(label string) *Option {
	o.Label = label
	return o
}

func (o *Option) AddOptions() string {
	return o.Label
}

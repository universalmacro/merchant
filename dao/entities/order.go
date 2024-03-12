package entities

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/snowflake"
	"gorm.io/gorm"
)

type Table struct {
	SpaceAsset
	Label string
}

func (a *Table) BeforeCreate(tx *gorm.DB) (err error) {
	a.Model.ID = snowflake.NewIdGenertor(0).Uint()
	return err
}

type Food struct {
	SpaceAsset
	Name        string `gorm:"type:varchar(100)"`
	Description string `gorm:"type:varchar(200)"`
	Price       int64
	FixedOffset *int64
	Image       string `gorm:"type:varchar(256)"`
	Categories  dao.StringArray
	Attributes  Attributes
	Status      string `gorm:"type:varchar(64)"`
	Printers    dao.UintArray
}

func (a *Food) BeforeCreate(tx *gorm.DB) (err error) {
	a.Model.ID = snowflake.NewIdGenertor(0).Uint()
	return err
}

func (Food) GormDataType() string {
	return "JSON"
}

func (s *Food) Scan(value any) error {
	return json.Unmarshal(value.([]byte), s)
}

func (s Food) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	return b, err
}

type Option struct {
	Label string `json:"label"`
	Extra int64  `json:"extra"`
}

type Options []Option

func (Options) GormDataType() string {
	return "JSON"
}

func (s *Options) Scan(value any) error {
	return json.Unmarshal(value.([]byte), s)
}

func (s Options) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	return b, err
}

type Attribute struct {
	Label   string   `json:"label"`
	Options []Option `json:"options"`
}

func (Attribute) GormDataType() string {
	return "JSON"
}

func (s *Attribute) Scan(value any) error {
	return json.Unmarshal(value.([]byte), s)
}

func (s Attribute) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	return b, err
}

type Attributes []Attribute

func (as Attributes) GetOption(left, right string) *Option {
	for _, a := range as {
		if left == a.Label {
			for _, option := range a.Options {
				if right == option.Label {
					return &option
				}
			}
		}
	}
	return nil
}

func (Attributes) GormDataType() string {
	return "JSON"
}

func (s *Attributes) Scan(value any) error {
	return json.Unmarshal(value.([]byte), s)
}

func (s Attributes) Value() (driver.Value, error) {
	return json.Marshal(s)
}

type FoodSpec struct {
	Food `json:"food"`
	Spec Spec `json:"spec"`
}

func (f *FoodSpec) Equals(from FoodSpec) bool {
	return true
}

func (f *FoodSpec) Scan(value any) error {
	return json.Unmarshal(value.([]byte), f)
}

func (f FoodSpec) Value() (driver.Value, error) {
	return json.Marshal(f)
}

func (f *FoodSpec) SetFood(food Food) *FoodSpec {
	f.Food = food
	return f
}

func (f *FoodSpec) SetSpec(spec Spec) *FoodSpec {
	f.Spec = spec
	return f
}

func (f *FoodSpec) Price() int64 {
	var total int64
	for _, spec := range f.Spec {
		option := f.Attributes.GetOption(spec.Attribute, spec.Optioned)
		if option != nil {
			total += option.Extra
		}
	}
	f.Food.Price += total
	return total
}

type Spec []SpecItem

type SpecItem struct {
	Attribute string `json:"attribute"`
	Optioned  string `json:"optioned"`
}

func (f Spec) Equals(from Spec) bool {
	selfMap := make(map[string]string)
	fromMap := make(map[string]string)
	for _, spec := range f {
		selfMap[spec.Attribute] = spec.Optioned
	}
	for _, spec := range from {
		fromMap[spec.Attribute] = spec.Optioned
	}
	for k, v := range selfMap {
		if fromMap[k] != v {
			return false
		}
	}
	return true
}

type Order struct {
	SpaceAsset
	PickUpCode int64
	TableLabel *string
	Foods      FoodSpces
	Status     string `gorm:"index;type:varchar(64)"`
	BillId     uint   `gorm:"index"`
	CreatorID  *int64
}

func (a *Order) BeforeCreate(tx *gorm.DB) (err error) {
	a.Model.ID = snowflake.NewIdGenertor(0).Uint()
	return err
}

type FoodSpces []FoodSpec

func (FoodSpces) GormDataType() string {
	return "JSON"
}

func (f *FoodSpces) Scan(value any) error {
	return json.Unmarshal(value.([]byte), f)
}

func (f FoodSpces) Value() (driver.Value, error) {
	return json.Marshal(f)
}

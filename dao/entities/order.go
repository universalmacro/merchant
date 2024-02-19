package entities

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/data"
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

func (as Attributes) GetOption(left, right string) (data.Pair[string, string], error) {
	for _, a := range as {
		if left == a.Label {
			for _, option := range a.Options {
				if right == option.Label {
					return data.Pair[string, string]{L: left, R: right}, nil
				}
			}
		}
	}
	return data.Pair[string, string]{}, errors.New("not found")
}

func (Attributes) GormDataType() string {
	return "JSON"
}

func (s *Attributes) Scan(value any) error {
	return json.Unmarshal(value.([]byte), s)
}

func (s Attributes) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	return b, err
}

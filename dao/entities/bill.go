package entities

import (
	"github.com/universalmacro/common/snowflake"
	"gorm.io/gorm"
)

type Bill struct {
	gorm.Model
	CashierID uint `json:"cashierId"`
	Amount    uint `json:"amount"`
}

func (a *Bill) BeforeCreate(tx *gorm.DB) (err error) {
	a.Model.ID = snowflake.NewIdGenertor(0).Uint()
	return err
}

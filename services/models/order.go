package models

import (
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/dao/repositories"
)

type Table struct {
	*entities.Table
}

func (t *Table) ID() uint {
	return t.Table.ID
}

func (t *Table) Label() string {
	return t.Table.Label
}

func (t *Table) SetLabel(label string) *Table {
	t.Table.Label = label
	return t
}

func (t *Table) Submit() *Table {
	repo := repositories.GetTableRepository()
	repo.Save(t.Table)
	return t
}

func (t *Table) Delete() {
	repo := repositories.GetTableRepository()
	repo.Delete(t.Table)
}

func (t *Table) SpaceID() uint {
	return t.Table.SpaceID
}

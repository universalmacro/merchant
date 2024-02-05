package models

import (
	"github.com/universalmacro/common/utils"
	"github.com/universalmacro/merchant/dao/entities"
)

type Space struct {
	*entities.Space
}

func (s *Space) ID() uint {
	return s.Space.ID
}

func (s *Space) StringID() string {
	return utils.UintToString(s.ID())
}

package service

import (
	"github.com/LeftUnion/theName/service/TheName/dtos"
	"github.com/LeftUnion/theName/service/TheName/internal/autogen"
)

type IService interface {
	// Human.
	GetHuman(id int) (*dtos.ResponseRichHuman, error)
	DeleteHuman(id int) error
	UpdateHuman(human autogen.UpdateHumanJSONBody) error

	// Humans.
	GetHumans(params autogen.GetHumansParams) ([]dtos.ResponseRichHuman, int, error)
	AddHumans(poorHumans autogen.AddHumansJSONBody) error
	DeleteHumans(ids []int) error
	UpdateHumans(humans autogen.UpdateHumansJSONBody) error
}

package reeychModelId

import (
	"github.com/lithammer/shortuuid"
)

type ReeychId struct {
	id    string
	model string
}

type ID interface {
	Gen(string) *ReeychId
	Regen()
}

func Gen(modelName string) *ReeychId {
	return &ReeychId{id: shortuuid.New(), model: modelName}
}

func (rId *ReeychId) Regen() {
	rId.id = shortuuid.New()
}

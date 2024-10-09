package uuid

import (
	"github.com/google/uuid"
)

type WithUUIDGeneratorable interface {
	SetUUIDGen(uuidFunc Generatorable)
	GetUUIDGen() Generatorable
}

type WithUUIDGenerator struct {
	uuidFunc Generatorable
}

var generator = defaultGenerator{}

func (w *WithUUIDGenerator) SetUUIDGen(uuidFunc Generatorable) {
	if uuidFunc == nil {
		uuidFunc = generator
	}

	w.uuidFunc = uuidFunc
}

func (w *WithUUIDGenerator) GetUUIDGen() Generatorable {
	if w.uuidFunc == nil {
		w.uuidFunc = generator
	}

	return w.uuidFunc
}

func (w *WithUUIDGenerator) UUID() uuid.UUID {
	if w.uuidFunc == nil {
		w.uuidFunc = generator
	}

	return w.uuidFunc.UUID()
}

package uuid

import (
	"github.com/google/uuid"
)

//go:generate mockgen -source=generator.go -destination=./generator_mock.go -package=uuid -mock_names Generatorable=Mock
type Generatorable interface {
	UUID() uuid.UUID
}

type defaultGenerator struct{}

func (d defaultGenerator) UUID() uuid.UUID {
	id, _ := uuid.NewV7()

	return id
}

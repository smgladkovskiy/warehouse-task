package entities

import (
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/now"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/uuid"
)

type Option[T any] func(options T) error

func WithNowFunc[T now.WithNowGeneratorable](nowFunc now.Generatorable) Option[T] {
	return func(options T) error {
		options.SetNowGen(nowFunc)

		return nil
	}
}

func WithUUIDFunc[T uuid.WithUUIDGeneratorable](uuidFunc uuid.Generatorable) Option[T] {
	return func(options T) error {
		options.SetUUIDGen(uuidFunc)

		return nil
	}
}

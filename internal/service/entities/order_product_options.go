package entities

import vObject "github.com/smgladkovskiy/warehouse-task/internal/service/entities/value_objects"

func WithOrderProductPrice(price vObject.Price) func(*OrderProduct) error {
	return func(op *OrderProduct) error {
		op.Price = price

		return nil
	}
}

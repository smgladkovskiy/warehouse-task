package upsertorderproduct

import "github.com/smgladkovskiy/warehouse-task/internal/service/entities"

type Command struct {
	orderProduct *entities.OrderProduct
}

func NewCommandUnsafe(orderProduct *entities.OrderProduct) Command {
	return Command{orderProduct: orderProduct}
}

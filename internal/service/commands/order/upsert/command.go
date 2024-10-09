package upsertorder

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/smgladkovskiy/warehouse-task/internal/service/entities"
)

type Command struct {
	order *entities.Order
}

func NewCommandUnsafe(order *entities.Order) Command {
	return Command{order: order}
}

func NewCommand(userID uuid.UUID, opts ...entities.Option[*entities.Order]) (*Command, error) {
	cart, err := entities.NewOrder(userID, opts...)
	if err != nil {
		return nil, fmt.Errorf("[NewCommand - entities.NewOrder error]: %w", err)
	}

	return &Command{order: cart}, nil
}

func (c Command) GetOrder() *entities.Order {
	return c.order
}

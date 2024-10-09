package orders

import (
	"context"

	"github.com/smgladkovskiy/warehouse-task/internal/service/entities"
	queryOptions "github.com/smgladkovskiy/warehouse-task/internal/service/entities/query_options"
)

func (r *Repository) GetOrder(ctx context.Context, qos queryOptions.OrderQueryOptionable) (*entities.Order, error) {
	//TODO implement me
	panic("implement me")
}

package getorder

import (
	"context"

	"github.com/smgladkovskiy/warehouse-task/internal/service/entities"
	queryOptions "github.com/smgladkovskiy/warehouse-task/internal/service/entities/query_options"
)

//go:generate mockgen -source=handler.go -destination=order_getter_mock.go -package=getorder -mock_names OrderGetter=GetOrderMock
type OrderGetter interface {
	GetOrder(ctx context.Context, qos queryOptions.OrderQueryOptionable) (*entities.Order, error)
}

type QueryHandler struct {
	repo OrderGetter
}

func NewQueryHandler(repo OrderGetter) *QueryHandler {
	if repo == nil {
		panic("OrderGetter repo is nil")
	}

	return &QueryHandler{repo: repo}
}

func (h *QueryHandler) Handle(ctx context.Context, q Query) (*entities.Order, error) {
	return h.repo.GetOrder(ctx, queryOptions.NewOrderQueryOptions(q.qos...))
}

package getproduct

import (
	"context"

	"github.com/smgladkovskiy/warehouse-task/internal/service/entities"
	queryOptions "github.com/smgladkovskiy/warehouse-task/internal/service/entities/query_options"
)

//go:generate mockgen -source=handler.go -destination=product_getter_mock.go -package=getproduct -mock_names ProductGetter=GetProductMock
type ProductGetter interface {
	GetProduct(ctx context.Context, qos queryOptions.ProductQueryOptionable) (*entities.Product, error)
}

type QueryHandler struct {
	repo ProductGetter
}

func NewQueryHandler(repo ProductGetter) *QueryHandler {
	if repo == nil {
		panic("ProductGetter repo is nil")
	}

	return &QueryHandler{repo: repo}
}

func (h *QueryHandler) Handle(ctx context.Context, q Query) (*entities.Product, error) {
	return h.repo.GetProduct(ctx, queryOptions.NewProductQueryOptions(q.qos...))
}

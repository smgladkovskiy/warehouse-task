package getstocks

import (
	"context"

	"github.com/smgladkovskiy/warehouse-task/internal/service/entities"
	queryOptions "github.com/smgladkovskiy/warehouse-task/internal/service/entities/query_options"
)

//go:generate mockgen -source=handler.go -destination=stocks_getter_mock.go -package=getstocks -mock_names StocksGetter=GetStocksMock
type StocksGetter interface {
	GetStocks(ctx context.Context, qos queryOptions.StockQueryOptionable) (entities.Stocks, error)
}

type QueryHandler struct {
	repo StocksGetter
}

func NewQueryHandler(repo StocksGetter) *QueryHandler {
	if repo == nil {
		panic("StocksGetter repo is nil")
	}

	return &QueryHandler{repo: repo}
}

func (h *QueryHandler) Handle(ctx context.Context, q Query) (entities.Stocks, error) {
	return h.repo.GetStocks(ctx, queryOptions.NewStockQueryOptions(q.qos...))
}

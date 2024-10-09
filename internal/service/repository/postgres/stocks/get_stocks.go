package stocks

import (
	"context"

	"github.com/smgladkovskiy/warehouse-task/internal/service/entities"
	queryOptions "github.com/smgladkovskiy/warehouse-task/internal/service/entities/query_options"
)

func (r *Repository) GetStocks(ctx context.Context, qos queryOptions.StockQueryOptionable) (entities.Stocks, error) {
	//TODO implement me
	panic("implement me")
}

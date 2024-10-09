package getstocks

import (
	queryOptions "github.com/smgladkovskiy/warehouse-task/internal/service/entities/query_options"
	vObject "github.com/smgladkovskiy/warehouse-task/internal/service/entities/value_objects"
)

type Query struct {
	qos []queryOptions.QueryOption[*queryOptions.StockQueryOptions]
}

func NewQueryByProductIDUnsafe(productID vObject.ProductID) Query {
	return Query{
		qos: []queryOptions.QueryOption[*queryOptions.StockQueryOptions]{queryOptions.WithStockProductID(productID)},
	}
}

package queryoptions

import vObject "github.com/smgladkovskiy/warehouse-task/internal/service/entities/value_objects"

type StockQueryOptionable interface {
	QueryOptionable
	MetaQueryOptionable

	ForProductID() *vObject.ProductID
}

type StockQueryOptions struct {
	BasicQueryOptions
	MetaQueryOptions

	productID vObject.ProductID
}

func (p StockQueryOptions) ForProductID() *vObject.ProductID {
	return &p.productID
}

type StockQueryOption func(options *StockQueryOptions)

var _ StockQueryOptionable = (*StockQueryOptions)(nil)

func NewStockQueryOptions(queryOption ...QueryOption[*StockQueryOptions]) *StockQueryOptions {
	qos := StockQueryOptions{
		BasicQueryOptions: *NewBasicQueryOptions(),
		MetaQueryOptions:  *NewMetaQueryOptions(),
	}

	for _, opt := range queryOption {
		opt(&qos)
	}

	return &qos
}

func WithStockProductID(productID vObject.ProductID) QueryOption[*StockQueryOptions] {
	return func(options *StockQueryOptions) {
		options.productID = productID
	}
}

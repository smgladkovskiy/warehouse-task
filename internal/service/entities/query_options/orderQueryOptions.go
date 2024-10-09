package queryoptions

import vObject "github.com/smgladkovskiy/warehouse-task/internal/service/entities/value_objects"

type OrderQueryOptionable interface {
	QueryOptionable
	MetaQueryOptionable

	ForOrderID() *vObject.OrderID
}

type OrderQueryOptions struct {
	BasicQueryOptions
	MetaQueryOptions

	orderID vObject.OrderID
}

func (p OrderQueryOptions) ForOrderID() *vObject.OrderID {
	return &p.orderID
}

type OrderQueryOption func(options *OrderQueryOptions)

var _ OrderQueryOptionable = (*OrderQueryOptions)(nil)

func NewOrderQueryOptions(queryOption ...QueryOption[*OrderQueryOptions]) *OrderQueryOptions {
	qos := OrderQueryOptions{
		BasicQueryOptions: *NewBasicQueryOptions(),
		MetaQueryOptions:  *NewMetaQueryOptions(),
	}

	for _, opt := range queryOption {
		opt(&qos)
	}

	return &qos
}

func WithOrderID(orderID vObject.OrderID) QueryOption[*OrderQueryOptions] {
	return func(options *OrderQueryOptions) {
		options.orderID = orderID
	}
}

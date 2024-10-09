package queryoptions

import vObject "github.com/smgladkovskiy/warehouse-task/internal/service/entities/value_objects"

type ProductQueryOptionable interface {
	QueryOptionable
	MetaQueryOptionable

	ForProductID() *vObject.ProductID
}

type ProductQueryOptions struct {
	BasicQueryOptions
	MetaQueryOptions

	productID vObject.ProductID
}

func (p ProductQueryOptions) ForProductID() *vObject.ProductID {
	return &p.productID
}

type ProductQueryOption func(options *ProductQueryOptions)

var _ ProductQueryOptionable = (*ProductQueryOptions)(nil)

func NewProductQueryOptions(queryOption ...QueryOption[*ProductQueryOptions]) *ProductQueryOptions {
	qos := ProductQueryOptions{
		BasicQueryOptions: *NewBasicQueryOptions(),
		MetaQueryOptions:  *NewMetaQueryOptions(),
	}

	for _, opt := range queryOption {
		opt(&qos)
	}

	return &qos
}

func WithProductID(productID vObject.ProductID) QueryOption[*ProductQueryOptions] {
	return func(options *ProductQueryOptions) {
		options.productID = productID
	}
}

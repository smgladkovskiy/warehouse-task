package getproduct

import (
	"fmt"

	"github.com/google/uuid"

	queryOptions "github.com/smgladkovskiy/warehouse-task/internal/service/entities/query_options"
	vObject "github.com/smgladkovskiy/warehouse-task/internal/service/entities/value_objects"
)

type Query struct {
	qos []queryOptions.QueryOption[*queryOptions.ProductQueryOptions]
}

func NewQueryByID(productUUID uuid.UUID) (*Query, error) {
	productID, err := vObject.NewProductIDFromUUID(productUUID)
	if err != nil {
		return nil, fmt.Errorf("[NewQueryByIDForUpdate - vObject.NewProductIDFromUUID error]: %w", err)
	}

	return &Query{
		qos: []queryOptions.QueryOption[*queryOptions.ProductQueryOptions]{queryOptions.WithProductID(productID)},
	}, nil
}

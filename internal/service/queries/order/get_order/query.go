package getorder

import (
	"github.com/google/uuid"

	queryOptions "github.com/smgladkovskiy/warehouse-task/internal/service/entities/query_options"
	vObject "github.com/smgladkovskiy/warehouse-task/internal/service/entities/value_objects"
)

type Query struct {
	qos []queryOptions.QueryOption[*queryOptions.OrderQueryOptions]
}

func NewQueryForUpdate(orderUUID uuid.UUID) (*Query, error) {
	orderID, err := vObject.NewOrderIDFromUUID(orderUUID)
	if err != nil {
		return nil, err
	}

	return &Query{
		qos: []queryOptions.QueryOption[*queryOptions.OrderQueryOptions]{
			queryOptions.WithOrderID(orderID),
			queryOptions.WithForUpdate[*queryOptions.OrderQueryOptions](),
		},
	}, nil
}

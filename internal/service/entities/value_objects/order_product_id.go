package valueobjects

import (
	"fmt"

	"github.com/google/uuid"
)

type OrderProductID struct {
	withUUIDer
}

func NewOrderProductIDFromUUID(id uuid.UUID) (OrderProductID, error) {
	if id == uuid.Nil {
		return OrderProductID{}, fmt.Errorf("orderProduct %w", ErrEmptyID)
	}

	return NewOrderProductIDFromUUIDUnsafe(id), nil
}

func NewOrderProductIDFromUUIDUnsafe(id uuid.UUID) OrderProductID {
	userID := OrderProductID{}
	userID.SetFromUUID(id)

	return userID
}

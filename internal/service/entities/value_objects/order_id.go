package valueobjects

import (
	"fmt"

	"github.com/google/uuid"
)

type OrderID struct {
	withUUIDer
}

func NewOrderIDFromUUID(id uuid.UUID) (OrderID, error) {
	if id == uuid.Nil {
		return OrderID{}, fmt.Errorf("order %w", ErrEmptyID)
	}

	return NewOrderIDFromUUIDUnsafe(id), nil
}

func NewOrderIDFromUUIDUnsafe(id uuid.UUID) OrderID {
	userID := OrderID{}
	userID.SetFromUUID(id)

	return userID
}

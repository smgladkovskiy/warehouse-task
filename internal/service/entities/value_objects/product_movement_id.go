package valueobjects

import (
	"fmt"

	"github.com/google/uuid"
)

type ProductMovementID struct {
	withUUIDer
}

func NewProductMovementIDFromUUID(id uuid.UUID) (ProductMovementID, error) {
	if id == uuid.Nil {
		return ProductMovementID{}, fmt.Errorf("product movement %w", ErrEmptyID)
	}

	return NewProductMovementIDFromUUIDUnsafe(id), nil
}

func NewProductMovementIDFromUUIDUnsafe(id uuid.UUID) ProductMovementID {
	userID := ProductMovementID{}
	userID.SetFromUUID(id)

	return userID
}

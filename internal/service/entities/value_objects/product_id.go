package valueobjects

import (
	"fmt"

	"github.com/google/uuid"
)

type ProductID struct {
	withUUIDer
}

func NewProductIDFromUUID(id uuid.UUID) (ProductID, error) {
	if id == uuid.Nil {
		return ProductID{}, fmt.Errorf("product %w", ErrEmptyID)
	}

	return NewProductIDFromUUIDUnsafe(id), nil
}

func NewProductIDFromUUIDUnsafe(id uuid.UUID) ProductID {
	userID := ProductID{}
	userID.SetFromUUID(id)

	return userID
}

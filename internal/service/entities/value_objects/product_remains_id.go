package valueobjects

import (
	"fmt"

	"github.com/google/uuid"
)

type ProductRemainsID struct {
	withUUIDer
}

func NewProductRemainsIDFromUUID(id uuid.UUID) (ProductRemainsID, error) {
	if id == uuid.Nil {
		return ProductRemainsID{}, fmt.Errorf("productRemains %w", ErrEmptyID)
	}

	return NewProductRemainsIDFromUUIDUnsafe(id), nil
}

func NewProductRemainsIDFromUUIDUnsafe(id uuid.UUID) ProductRemainsID {
	userID := ProductRemainsID{}
	userID.SetFromUUID(id)

	return userID
}

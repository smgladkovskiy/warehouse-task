package valueobjects

import (
	"fmt"

	"github.com/google/uuid"
)

type WarehouseID struct {
	withUUIDer
}

func NewWarehouseIDFromUUID(id uuid.UUID) (WarehouseID, error) {
	if id == uuid.Nil {
		return WarehouseID{}, fmt.Errorf("warehouse %w", ErrEmptyID)
	}

	return NewWarehouseIDFromUUIDUnsafe(id), nil
}

func NewWarehouseIDFromUUIDUnsafe(id uuid.UUID) WarehouseID {
	userID := WarehouseID{}
	userID.SetFromUUID(id)

	return userID
}

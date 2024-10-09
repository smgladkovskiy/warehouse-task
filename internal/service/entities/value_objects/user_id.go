package valueobjects

import (
	"fmt"

	"github.com/google/uuid"
)

type UserID struct {
	withUUIDer
}

func NewUserIDFromUUID(id uuid.UUID) (UserID, error) {
	if id == uuid.Nil {
		return UserID{}, fmt.Errorf("user %w", ErrEmptyID)
	}

	return NewUserIDFromUUIDUnsafe(id), nil
}

func NewUserIDFromUUIDUnsafe(id uuid.UUID) UserID {
	userID := UserID{}
	userID.SetFromUUID(id)

	return userID
}

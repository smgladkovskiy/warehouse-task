package valueobjects

import "errors"

type MaritalStatus string

const (
	MaritalStatusUnknown  MaritalStatus = "unknown"
	MaritalStatusSingle   MaritalStatus = "single"
	MaritalStatusMarried  MaritalStatus = "married"
	MaritalStatusDivorced MaritalStatus = "divorced"
	MaritalStatusWidowed  MaritalStatus = "widowed"
)

var availableMaritalStatuses = map[MaritalStatus]struct{}{
	MaritalStatusSingle:   {},
	MaritalStatusMarried:  {},
	MaritalStatusDivorced: {},
	MaritalStatusWidowed:  {},
}

var (
	ErrEmptyMaritalStatus   = errors.New("empty marital status")
	ErrUnknownMaritalStatus = errors.New("unknown marital status")
)

func NewMaritalStatus(status string) (MaritalStatus, error) {
	if status == "" {
		return MaritalStatusUnknown, ErrEmptyMaritalStatus
	}

	ms := NewMaritalStatusUnsafe(status)

	if _, ok := availableMaritalStatuses[ms]; !ok {
		return MaritalStatusUnknown, ErrUnknownMaritalStatus
	}

	return ms, nil
}

func NewMaritalStatusUnsafe(status string) MaritalStatus {
	return MaritalStatus(status)
}

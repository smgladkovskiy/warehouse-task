package valueobjects

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
)

type withUUIDer struct {
	id uuid.UUID
}

func NilUUIDer() withUUIDer {
	return withUUIDer{id: uuid.Nil}
}

func NewIDFromStr(uuidStr string) (withUUIDer, error) {
	if uuidStr == "" {
		return NilUUIDer(), nil
	}

	id, err := uuid.Parse(uuidStr)
	if err != nil {
		return NilUUIDer(), fmt.Errorf("new ID from string %w: %s", ErrParseID, err.Error())
	}

	return withUUIDer{id: id}, nil
}

func (id withUUIDer) Clone() IDInterface {
	newDID := id

	return &newDID
}

func (id withUUIDer) Name() string {
	return "withUUIDer"
}

func (id *withUUIDer) SetEmpty() {
	*id = withUUIDer{}
}

func (id *withUUIDer) SetFromString(sv driver.Value) error {
	var err error

	*id, err = NewIDFromStr(fmt.Sprintf("%s", sv))
	if err != nil {
		return err
	}

	return nil
}

func (id *withUUIDer) SetFromUUID(value uuid.UUID) {
	*id = withUUIDer{id: value}
}

func (id withUUIDer) IsNil() bool {
	return id.UUID() == uuid.Nil
}

func (id withUUIDer) String() string {
	return id.UUID().String()
}

func (id withUUIDer) Bytes() []byte {
	return idBytes(&id)
}

func (id withUUIDer) UUID() uuid.UUID {
	return id.id
}

func (id withUUIDer) MarshalJSON() ([]byte, error) {
	return idMarshalJSON(&id)
}

func (id *withUUIDer) UnmarshalJSON(data []byte) error {
	return idUnmarshalJSON(data, id)
}

func (id *withUUIDer) Scan(value interface{}) error {
	return idScan(id, value)
}

func (id withUUIDer) Value() (driver.Value, error) {
	return getValue(&id)
}

func (id withUUIDer) EqualsAnother(anotherID IDInterface) bool {
	return id.String() == anotherID.String()
}

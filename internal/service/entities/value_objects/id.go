package valueobjects

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// IDInterface - это интерфейс для идентификатора модели, представленной в формате UUID.
// Унифицирует последующую работу с идентификатором и предоставляет возможность проводить тесты.
type IDInterface interface {
	Clone() IDInterface
	Name() string
	SetEmpty()
	SetFromString(sv driver.Value) error
	SetFromUUID(value uuid.UUID)
	IsNil() bool
	UUID() uuid.UUID
	String() string
	Bytes() []byte
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
	Scan(value interface{}) error
	Value() (driver.Value, error)
	EqualsAnother(id IDInterface) bool
}

const NilUUIDStr = "00000000-0000-0000-0000-000000000000"

var (
	ErrScanError = errors.New("scan error for id")
	ErrEmptyID   = errors.New("id has nil value")
	ErrParseID   = errors.New("id parsing error")
)

func getValue(id IDInterface) (driver.Value, error) {
	if id.IsNil() {
		return nil, nil
	}

	return id.String(), nil
}

func idBytes(id IDInterface) []byte {
	if id.IsNil() {
		bbNil, _ := uuid.Nil.MarshalText() //olint:errcheck // тут ошибки быть не может

		return bbNil
	}

	return []byte(id.String())
}

func idMarshalJSON(id IDInterface) ([]byte, error) {
	bb := &bytes.Buffer{}

	uuidBb, err := id.UUID().MarshalText()
	if err != nil {
		return nil, fmt.Errorf("%s.MarshalJSON uuid.MarshalText error: %w", id.Name(), err)
	}

	bb.WriteString(`"`)

	for _, b := range uuidBb {
		bb.WriteByte(b)
	}

	bb.WriteString(`"`)

	return bb.Bytes(), nil
}

func idScan(id IDInterface, value interface{}) error {
	// if value is nil, false
	if value == nil {
		// set the value of the pointer to empty UUID
		id.SetEmpty()

		return nil
	}

	stringValue, err := driver.String.ConvertValue(value)
	if err != nil {
		return fmt.Errorf("%w - %s: %s", ErrScanError, id.Name(), err.Error())
	}

	if stringValue == "{}" {
		return fmt.Errorf("%w - %s", ErrScanError, id.Name())
	}

	return id.SetFromString(stringValue)
}

func idUnmarshalJSON(data []byte, id IDInterface) error {
	var idStr string

	if data == nil {
		id.SetEmpty()

		return nil
	}

	if err := json.Unmarshal(data, &idStr); err != nil {
		return fmt.Errorf("%s.UnmarshalJSON error: %w", id.Name(), err)
	}

	return id.SetFromString(idStr)
}

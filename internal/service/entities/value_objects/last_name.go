package valueobjects

import "errors"

type LastName string

const LastNameEmpty LastName = ""

var ErrEmptyLastName = errors.New("empty last name")

func NewLastName(name string) (LastName, error) {
	if name == "" {
		return LastNameEmpty, ErrEmptyLastName
	}

	return NewLastNameUnsafe(name), nil
}

func NewLastNameUnsafe(name string) LastName {
	return LastName(name)
}

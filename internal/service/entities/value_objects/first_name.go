package valueobjects

import "errors"

type FirstName string

const FirstNameEmpty FirstName = ""

var ErrEmptyFirstName = errors.New("empty first name")

func NewFirstName(name string) (FirstName, error) {
	if name == "" {
		return FirstNameEmpty, ErrEmptyFirstName
	}

	return NewFirstNameUnsafe(name), nil
}

func NewFirstNameUnsafe(name string) FirstName {
	return FirstName(name)
}

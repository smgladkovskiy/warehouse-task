package valueobjects

import (
	"errors"
	"fmt"
	"time"
)

type Birthdate time.Time

const (
	minAge = 18
)

var AgeUnknown = Birthdate(time.Time{})

var ErrAgeIsTooLow = errors.New("age must be greater than")

func NewBirthDate(birthDate time.Time) (Birthdate, error) {
	bd := NewAgeUnsafe(birthDate)

	if bd.Age() < minAge {
		return AgeUnknown, fmt.Errorf("%w %d", ErrAgeIsTooLow, minAge)
	}

	return NewAgeUnsafe(birthDate), nil
}

func NewAgeUnsafe(birthDate time.Time) Birthdate {
	return Birthdate(time.Date(birthDate.Year(), birthDate.Month(), birthDate.Day(), 0, 0, 0, 0, time.UTC))
}

func (b Birthdate) Age() int {
	if b == AgeUnknown {
		return 0
	}

	return time.Now().Year() - b.Time().Year()
}

func (b Birthdate) Time() time.Time {
	return time.Time(b)
}

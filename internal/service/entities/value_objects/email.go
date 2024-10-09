package valueobjects

import (
	"errors"
	"net/mail"
)

type Email string

const EmailEmpty Email = ""

var ErrEmptyEmail = errors.New("email is empty")

func NewEmail(email string) (Email, error) {
	if email == "" {
		return EmailEmpty, ErrEmptyEmail
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return EmailEmpty, err
	}

	return NewEmailUnsafe(email), nil
}

func NewEmailUnsafe(email string) Email {
	return Email(email)
}

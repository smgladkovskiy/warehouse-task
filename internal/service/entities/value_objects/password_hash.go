package valueobjects

import (
	"errors"
	"fmt"

	passCrypto "github.com/smgladkovskiy/warehouse-task/internal/pkg/pass_crypto"
)

type PasswordHash string

const (
	minPasswordLength              = 8
	PasswordHashEmpty PasswordHash = ""
)

var (
	ErrPasswordLen = errors.New("password length must be at least")
)

func NewPasswordHash(hasher passCrypto.PasswordHashable, password string) (PasswordHash, error) {
	if len(password) < minPasswordLength {
		return PasswordHashEmpty, fmt.Errorf("%w %d symbols but got only %d", ErrPasswordLen, minPasswordLength, len(password))
	}

	hash, err := hasher.HashAndSalt([]byte(password))
	if err != nil {
		return PasswordHashEmpty, err
	}

	return NewPasswordHashUnsafe(hash), nil
}

func NewPasswordHashUnsafe(hash string) PasswordHash {
	return PasswordHash(hash)
}

package passcrypto

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=hasher.go -destination=hasher_mock.go -package=passcrypto -mock_names PasswordHashable=PasswordHashMock
type PasswordHashable interface {
	HashAndSalt(pwd []byte) (string, error)
	ComparePasswords(hashedPwd string, plainPwd []byte) bool
}

type defaultPasswordHasher struct{}

var ErrPasswordHashing = errors.New("password hashing error")

// HashAndSalt Hashes a given string
func (d defaultPasswordHasher) HashAndSalt(pwd []byte) (string, error) {

	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrPasswordHashing, err.Error())
	}

	return string(hash), nil
}

func (d defaultPasswordHasher) ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return false
	}

	return true
}

package entities

import (
	"fmt"

	passcrypto "github.com/smgladkovskiy/warehouse-task/internal/pkg/pass_crypto"
	vObject "github.com/smgladkovskiy/warehouse-task/internal/service/entities/value_objects"
)

func WithUserPassword(password string) func(*User) error {
	return func(u *User) error {
		var err error

		u.PasswordHash, err = vObject.NewPasswordHash(u.GetHasher(), password)
		if err != nil {
			return fmt.Errorf("[WithUserPassword] %w", err)
		}

		return nil
	}
}

func WithUserPasswordHasher(hasher passcrypto.PasswordHashable) func(*User) error {
	return func(u *User) error {
		u.SetHasher(hasher)

		return nil
	}
}

package userregistration

import (
	"fmt"

	passcrypto "github.com/smgladkovskiy/warehouse-task/internal/pkg/pass_crypto"
	createUser "github.com/smgladkovskiy/warehouse-task/internal/service/commands/user/create"
	getUserByEmail "github.com/smgladkovskiy/warehouse-task/internal/service/queries/user/get_by_email"
	usecase "github.com/smgladkovskiy/warehouse-task/internal/service/usecases"
)

func WithGetUserByEmailQuery(handler *getUserByEmail.QueryHandler) usecase.Configuration[*UseCase] {
	return func(uc *UseCase) error {
		if handler == nil {
			return fmt.Errorf("%w %s", usecase.ErrEmptyStructParam, "getUserByEmail")
		}

		uc.getUserQuery = handler

		return nil
	}
}

func WithCreateUserCommand(handler *createUser.CommandHandler) usecase.Configuration[*UseCase] {
	return func(uc *UseCase) error {
		if handler == nil {
			return fmt.Errorf("%w %s", usecase.ErrEmptyStructParam, "createUser")
		}

		uc.createUserCmd = handler

		return nil
	}
}

func WithPasswordHasher(hasher passcrypto.PasswordHashable) usecase.Configuration[*UseCase] {
	return func(uc *UseCase) error {
		if hasher == nil {
			return fmt.Errorf("%w %s", usecase.ErrEmptyStructParam, "hasher")
		}

		uc.SetHasher(hasher)

		return nil
	}

}

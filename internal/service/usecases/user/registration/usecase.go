package userregistration

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/checker"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/log"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/now"
	passcrypto "github.com/smgladkovskiy/warehouse-task/internal/pkg/pass_crypto"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/uuid"
	createUser "github.com/smgladkovskiy/warehouse-task/internal/service/commands/user/create"
	"github.com/smgladkovskiy/warehouse-task/internal/service/entities"
	getUserByEmail "github.com/smgladkovskiy/warehouse-task/internal/service/queries/user/get_by_email"
	usecase "github.com/smgladkovskiy/warehouse-task/internal/service/usecases"
)

type UseCase struct {
	uuid.WithUUIDGenerator
	now.WithNowGenerator
	checker.WithCheck
	passcrypto.WithPasswordHasher
	log.WithLogger

	// Query handlers
	getUserQuery *getUserByEmail.QueryHandler

	// Command handlers
	createUserCmd *createUser.CommandHandler
}

//var _ handler.Authenticator = (*UseCase)(nil)

func NewUseCase(cfgs ...usecase.Configuration[*UseCase]) (*UseCase, error) {
	uc := &UseCase{}

	// Apply all Configurations passed in
	for _, cfg := range cfgs {
		if cfg == nil {
			return nil, checker.ErrInitError
		}

		err := cfg(uc)
		if err != nil {
			return nil, err
		}
	}

	if err := uc.Check(*uc); err != nil {
		return nil, err
	}

	return uc, nil
}

func (uc *UseCase) Run(ctx context.Context, req Requestable) (*entities.User, error) {
	l := uc.Logger().With(
		log.String("email", req.GetEmail()),
		log.String("firstName", req.GetFirstName()),
		log.String("lastName", req.GetLastName()),
		log.String("birthDate", req.GetBirthDate().String()),
		log.String("maritalStatus", req.GetMaritalStatus()),
		log.String("password", strings.Repeat("*", len(req.GetPassword()))),
	)

	l.Debug(ctx, "START usecase")

	// 1. проверить наличие пользователя по email
	query, err := getUserByEmail.NewQuery(req.GetEmail())
	if err != nil {
		l.Error(ctx, "STOP usecase! getUserByEmail.NewQuery error", log.Err(err))

		return nil, fmt.Errorf("[userRegistration - getUserByEmail.NewQuery error]: %w", err)
	}

	existedUser, err := uc.getUserQuery.Handle(ctx, *query)
	if err != nil && !errors.Is(err, entities.ErrUserRecNotFound) {
		l.Error(ctx, "STOP usecase! getUserQuery.Handle error", log.Err(err))

		return nil, fmt.Errorf("[userRegistration - getUserQuery.Handle error]: %w", err)
	}

	if existedUser != nil {
		l.Error(ctx, "STOP usecase! user already exists", log.Err(entities.ErrUserAlreadyExists))

		return nil, fmt.Errorf("[userRegistration - Run error]: %w", entities.ErrUserAlreadyExists)
	}

	// 2. сохранить (зарегистрировать) пользователя
	cmd, err := createUser.NewCommand(
		req.GetEmail(),
		req.GetFirstName(),
		req.GetLastName(),
		req.GetMaritalStatus(),
		req.GetBirthDate(),
		entities.WithUserPasswordHasher(uc.GetHasher()),
		entities.WithUserPassword(req.GetPassword()),
		entities.WithUUIDFunc[*entities.User](uc.GetUUIDGen()),
		entities.WithNowFunc[*entities.User](uc.GetNowGen()),
	)
	if err != nil {
		l.Error(ctx, "STOP usecase! createUser.NewCommand error", log.Err(err))

		return nil, fmt.Errorf("[userRegistration - createUser.NewCommand error]: %w", err)
	}

	if err = uc.createUserCmd.Handle(ctx, *cmd); err != nil {
		l.Error(ctx, "STOP usecase! createUserCmd.Handle error", log.Err(err))

		return nil, fmt.Errorf("[userRegistration - createUserCmd.Handle error]: %w", err)
	}

	l.Debug(ctx, "END usecase")

	return cmd.GetUser(), nil
}

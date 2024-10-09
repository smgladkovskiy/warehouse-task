package userregistration

import (
	"context"
	"strings"
	"testing"
	"time"

	baseUUID "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/checker"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/log"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/now"
	passcrypto "github.com/smgladkovskiy/warehouse-task/internal/pkg/pass_crypto"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/uuid"
	createUser "github.com/smgladkovskiy/warehouse-task/internal/service/commands/user/create"
	"github.com/smgladkovskiy/warehouse-task/internal/service/entities"
	vObject "github.com/smgladkovskiy/warehouse-task/internal/service/entities/value_objects"
	getUserByEmail "github.com/smgladkovskiy/warehouse-task/internal/service/queries/user/get_by_email"
	usecase "github.com/smgladkovskiy/warehouse-task/internal/service/usecases"
)

func TestNewUseCase(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	getUserByEmailMock := getUserByEmail.NewGetUserMock(ctrl)
	createUserMock := createUser.NewCreateUserMock(ctrl)

	cfgs := []usecase.Configuration[*UseCase]{
		WithGetUserByEmailQuery(getUserByEmail.NewQueryHandler(getUserByEmailMock)),
		WithCreateUserCommand(createUser.NewCommandHandler(createUserMock)),
	}

	uc, err := NewUseCase(cfgs...)
	require.NoError(t, err)
	require.NotEmpty(t, uc)

	uc, err = NewUseCase(func(upc *UseCase) error {
		return assert.AnError
	})
	require.ErrorIs(t, err, assert.AnError)
	require.Empty(t, uc)

	uc, err = NewUseCase()
	require.ErrorIs(t, err, checker.ErrInitError)
	require.Empty(t, uc)

	uc, err = NewUseCase(nil)
	require.ErrorIs(t, err, checker.ErrInitError)
	require.Empty(t, uc)
}

func TestUseCase_Run(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string
		in   testRequest
		exp  func(t *testing.T, in testRequest, loggerMock *log.LogMock, hasherMock *passcrypto.PasswordHashMock, getUserByEmailMock *getUserByEmail.GetUserMock, createUserMock *createUser.CreateUserMock) (*entities.User, error)
	}

	id := baseUUID.New()
	tn := time.Now().UTC()
	ctrl := gomock.NewController(t)
	nowFunc := now.NewMock(ctrl)
	uuidFunc := uuid.NewMock(ctrl)

	nowFunc.EXPECT().Now().AnyTimes().Return(tn)
	uuidFunc.EXPECT().UUID().AnyTimes().Return(id)

	tcs := []testCase{
		{
			name: "happy path",
			in: testRequest{
				email:         "some@email.com",
				firstName:     "first name",
				lastName:      "last name",
				birthdate:     time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				maritalStatus: "married",
				password:      "12345678",
			},
			exp: func(t *testing.T, in testRequest, loggerMock *log.LogMock, hasherMock *passcrypto.PasswordHashMock, getUserByEmailMock *getUserByEmail.GetUserMock, createUserMock *createUser.CreateUserMock) (*entities.User, error) {
				t.Helper()

				hasherMock.EXPECT().HashAndSalt([]byte(in.GetPassword())).Return("hashed_password", nil).AnyTimes()

				user, err := entities.NewUser(
					in.GetEmail(),
					in.GetFirstName(),
					in.GetLastName(),
					in.GetMaritalStatus(),
					in.GetBirthDate(),
					entities.WithUserPasswordHasher(hasherMock),
					entities.WithUserPassword(in.GetPassword()),
					entities.WithUUIDFunc[*entities.User](uuidFunc),
					entities.WithNowFunc[*entities.User](nowFunc),
				)
				require.NoError(t, err)

				getUserByEmailMock.EXPECT().GetByEmail(gomock.Any(), user.Email).Return(nil, entities.ErrUserRecNotFound)
				createUserMock.EXPECT().CreateUser(gomock.Any(), user).Return(nil)
				loggerMock.EXPECT().Debug(gomock.Any(), "END usecase")

				return user, nil
			},
		},
		{
			name: "create user handler error",
			in: testRequest{
				email:         "some@email.com",
				firstName:     "first name",
				lastName:      "last name",
				birthdate:     time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				maritalStatus: "married",
				password:      "12345678",
			},
			exp: func(t *testing.T, in testRequest, loggerMock *log.LogMock, hasherMock *passcrypto.PasswordHashMock, getUserByEmailMock *getUserByEmail.GetUserMock, createUserMock *createUser.CreateUserMock) (*entities.User, error) {
				t.Helper()

				hasherMock.EXPECT().HashAndSalt([]byte(in.GetPassword())).Return("hashed_password", nil).AnyTimes()

				user, err := entities.NewUser(
					in.GetEmail(),
					in.GetFirstName(),
					in.GetLastName(),
					in.GetMaritalStatus(),
					in.GetBirthDate(),
					entities.WithUserPasswordHasher(hasherMock),
					entities.WithUserPassword(in.GetPassword()),
					entities.WithUUIDFunc[*entities.User](uuidFunc),
					entities.WithNowFunc[*entities.User](nowFunc),
				)
				require.NoError(t, err)

				getUserByEmailMock.EXPECT().GetByEmail(gomock.Any(), user.Email).Return(nil, entities.ErrUserRecNotFound)
				createUserMock.EXPECT().CreateUser(gomock.Any(), user).Return(assert.AnError)
				loggerMock.EXPECT().Error(gomock.Any(), "STOP usecase! createUserCmd.Handle error", log.Err(assert.AnError))

				return nil, assert.AnError
			},
		},
		{
			name: "create user cmd error",
			in: testRequest{
				email:         "some@email.com",
				firstName:     "first name",
				lastName:      "last name",
				birthdate:     time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				maritalStatus: "married",
				password:      "123456",
			},
			exp: func(t *testing.T, in testRequest, loggerMock *log.LogMock, hasherMock *passcrypto.PasswordHashMock, getUserByEmailMock *getUserByEmail.GetUserMock, createUserMock *createUser.CreateUserMock) (*entities.User, error) {
				t.Helper()

				hasherMock.EXPECT().HashAndSalt([]byte(in.GetPassword())).Return("hashed_password", nil).AnyTimes()

				getUserByEmailMock.EXPECT().GetByEmail(gomock.Any(), vObject.NewEmailUnsafe(in.GetEmail())).Return(nil, entities.ErrUserRecNotFound)
				loggerMock.EXPECT().Error(gomock.Any(), "STOP usecase! createUser.NewCommand error", gomock.Any())

				return nil, vObject.ErrPasswordLen
			},
		},
		{
			name: "user exists error",
			in: testRequest{
				email:         "some@email.com",
				firstName:     "first name",
				lastName:      "last name",
				birthdate:     time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				maritalStatus: "married",
				password:      "12345678",
			},
			exp: func(t *testing.T, in testRequest, loggerMock *log.LogMock, hasherMock *passcrypto.PasswordHashMock, getUserByEmailMock *getUserByEmail.GetUserMock, createUserMock *createUser.CreateUserMock) (*entities.User, error) {
				t.Helper()

				hasherMock.EXPECT().HashAndSalt([]byte(in.GetPassword())).Return("hashed_password", nil).AnyTimes()

				user, err := entities.NewUser(
					in.GetEmail(),
					in.GetFirstName(),
					in.GetLastName(),
					in.GetMaritalStatus(),
					in.GetBirthDate(),
					entities.WithUserPasswordHasher(hasherMock),
					entities.WithUserPassword(in.GetPassword()),
					entities.WithUUIDFunc[*entities.User](uuidFunc),
					entities.WithNowFunc[*entities.User](nowFunc),
				)
				require.NoError(t, err)

				getUserByEmailMock.EXPECT().GetByEmail(gomock.Any(), vObject.NewEmailUnsafe(in.GetEmail())).Return(user, nil)
				loggerMock.EXPECT().Error(gomock.Any(), "STOP usecase! user already exists", gomock.Any())

				return nil, entities.ErrUserAlreadyExists
			},
		},
		{
			name: "get user handler error",
			in: testRequest{
				email:         "some@email.com",
				firstName:     "first name",
				lastName:      "last name",
				birthdate:     time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				maritalStatus: "married",
				password:      "12345678",
			},
			exp: func(t *testing.T, in testRequest, loggerMock *log.LogMock, hasherMock *passcrypto.PasswordHashMock, getUserByEmailMock *getUserByEmail.GetUserMock, createUserMock *createUser.CreateUserMock) (*entities.User, error) {
				t.Helper()

				getUserByEmailMock.EXPECT().GetByEmail(gomock.Any(), vObject.NewEmailUnsafe(in.GetEmail())).Return(nil, assert.AnError)
				loggerMock.EXPECT().Error(gomock.Any(), "STOP usecase! getUserQuery.Handle error", gomock.Any())

				return nil, assert.AnError
			},
		},
		{
			name: "get user query error",
			in: testRequest{
				email:         "",
				firstName:     "first name",
				lastName:      "last name",
				birthdate:     time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				maritalStatus: "married",
				password:      "12345678",
			},
			exp: func(t *testing.T, in testRequest, loggerMock *log.LogMock, hasherMock *passcrypto.PasswordHashMock, getUserByEmailMock *getUserByEmail.GetUserMock, createUserMock *createUser.CreateUserMock) (*entities.User, error) {
				t.Helper()

				loggerMock.EXPECT().Error(gomock.Any(), "STOP usecase! getUserByEmail.NewQuery error", gomock.Any())

				return nil, vObject.ErrEmptyEmail
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrlT := gomock.NewController(t)
			loggerMock := log.NewLogMock(ctrlT)
			passHasherMock := passcrypto.NewPasswordHashMock(ctrlT)
			getUserByEmailMock := getUserByEmail.NewGetUserMock(ctrlT)
			createUserMock := createUser.NewCreateUserMock(ctrlT)

			cfgs := []usecase.Configuration[*UseCase]{
				usecase.WithLogger[*UseCase](loggerMock),
				usecase.WithNowFunc[*UseCase](nowFunc),
				usecase.WithUUIDFunc[*UseCase](uuidFunc),
				WithPasswordHasher(passHasherMock),
				WithGetUserByEmailQuery(getUserByEmail.NewQueryHandler(getUserByEmailMock)),
				WithCreateUserCommand(createUser.NewCommandHandler(createUserMock)),
			}

			loggerMock.EXPECT().With(
				log.String("email", tc.in.GetEmail()),
				log.String("firstName", tc.in.GetFirstName()),
				log.String("lastName", tc.in.GetLastName()),
				log.String("birthDate", tc.in.GetBirthDate().String()),
				log.String("maritalStatus", tc.in.GetMaritalStatus()),
				log.String("password", strings.Repeat("*", len(tc.in.GetPassword()))),
			).Return(loggerMock)
			loggerMock.EXPECT().Debug(gomock.Any(), "START usecase")

			expOut, expErr := tc.exp(t, tc.in, loggerMock, passHasherMock, getUserByEmailMock, createUserMock)

			uc, err := NewUseCase(cfgs...)
			require.NoError(t, err)

			out, err := uc.Run(context.Background(), tc.in)

			assert.ErrorIs(t, err, expErr)
			assert.Equal(t, expOut, out)
		})
	}
}

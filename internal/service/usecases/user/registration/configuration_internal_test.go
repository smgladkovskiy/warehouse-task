package userregistration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/checker"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/log"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/now"
	passCrypto "github.com/smgladkovskiy/warehouse-task/internal/pkg/pass_crypto"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/uuid"
	createUser "github.com/smgladkovskiy/warehouse-task/internal/service/commands/user/create"
	getUserByEmail "github.com/smgladkovskiy/warehouse-task/internal/service/queries/user/get_by_email"
	usecase "github.com/smgladkovskiy/warehouse-task/internal/service/usecases"
)

func TestConfiguration(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	loggerMock := log.NewLogMock(ctrl)
	nowFunc := now.NewMock(ctrl)
	uuidFunc := uuid.NewMock(ctrl)
	hasherMock := passCrypto.NewPasswordHashMock(ctrl)
	getUserByEmailMock := getUserByEmail.NewGetUserMock(ctrl)
	createUserMock := createUser.NewCreateUserMock(ctrl)

	cfgs := []usecase.Configuration[*UseCase]{
		usecase.WithLogger[*UseCase](loggerMock),
		usecase.WithNowFunc[*UseCase](nowFunc),
		usecase.WithUUIDFunc[*UseCase](uuidFunc),
		WithPasswordHasher(hasherMock),
		WithGetUserByEmailQuery(getUserByEmail.NewQueryHandler(getUserByEmailMock)),
		WithCreateUserCommand(createUser.NewCommandHandler(createUserMock)),
	}

	f := WithGetUserByEmailQuery(nil)
	uc, err := NewUseCase(f)
	require.Error(t, err)
	assert.Empty(t, uc)

	f = WithPasswordHasher(nil)
	uc, err = NewUseCase(f)
	require.Error(t, err)
	assert.Empty(t, uc)

	f = WithCreateUserCommand(nil)
	uc, err = NewUseCase(f)
	require.Error(t, err)
	assert.Empty(t, uc)

	uc, err = NewUseCase(nil)
	require.ErrorIs(t, err, checker.ErrInitError)
	require.Empty(t, uc)

	uc, err = NewUseCase(cfgs...)
	require.NoError(t, err)
	assert.NotEmpty(t, uc)
}

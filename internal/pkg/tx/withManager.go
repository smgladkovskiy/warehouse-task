package tx

import (
	"context"
	"errors"

	"github.com/avito-tech/go-transaction-manager/trm"
	_ "go.uber.org/mock/mockgen/model"
)

//go:generate mockgen -destination=manager_mock.go -package=tx -mock_names Manager=TransactionManagerMock github.com/avito-tech/go-transaction-manager/trm Manager
type TransactionManager interface {
	SetTrxManager(trxManager trm.Manager) error
	TransactionDo(ctx context.Context, trx func(ctx context.Context) error) error
}

type WithTransactionManager struct {
	trxManager trm.Manager
}

var ErrNilTransactionManager = errors.New("transaction manager is nil")

func (m *WithTransactionManager) SetTrxManager(trxManager trm.Manager) error {
	if trxManager == nil {
		return ErrNilTransactionManager
	}

	m.trxManager = trxManager

	return nil
}

func (m *WithTransactionManager) TransactionDo(ctx context.Context, trx func(ctx context.Context) error) error {
	return m.trxManager.Do(ctx, trx)
}

package usecase

import (
	"github.com/avito-tech/go-transaction-manager/trm"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/log"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/now"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/tx"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/uuid"
)

// Configuration Функциональные опции для юзкейсов.
type Configuration[T any] func(uc T) error

// WithLogger Конфигурирует логгер юзкейса
func WithLogger[T log.WithLoggerable](l log.Logger) Configuration[T] {
	return func(uc T) error {
		uc.SetLogger(l)

		return nil
	}
}

// WithTransactionManager конфигурирует транзакционный менеджер юзкейса.
func WithTransactionManager[T tx.TransactionManager](trxManager trm.Manager) Configuration[T] {
	return func(uc T) error {
		return uc.SetTrxManager(trxManager)
	}
}

// WithNowFunc Конфигурирует генератор текущей метки времени, который может использоваться в юзкейсе.
func WithNowFunc[T now.WithNowGeneratorable](nowFunc now.Generatorable) Configuration[T] {
	return func(uc T) error {
		uc.SetNowGen(nowFunc)

		return nil
	}
}

// WithUUIDFunc Конфигурирует генератор UUID, который может использоваться в юзкейсе.
func WithUUIDFunc[T uuid.WithUUIDGeneratorable](uuidFunc uuid.Generatorable) Configuration[T] {
	return func(uc T) error {
		uc.SetUUIDGen(uuidFunc)

		return nil
	}
}

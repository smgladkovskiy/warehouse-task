package users

import (
	trmgorm "github.com/avito-tech/go-transaction-manager/gorm"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/db"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/now"
	trx "github.com/smgladkovskiy/warehouse-task/internal/pkg/tx"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/uuid"
	createUser "github.com/smgladkovskiy/warehouse-task/internal/service/commands/user/create"
	getUserByEmail "github.com/smgladkovskiy/warehouse-task/internal/service/queries/user/get_by_email"
)

type Repository struct {
	now.WithNowGenerator
	uuid.WithUUIDGenerator
	trx.WithTransactionDB
}

var (
	_ createUser.UserCreator    = (*Repository)(nil)
	_ getUserByEmail.UserGetter = (*Repository)(nil)
)

func NewRepository(db *db.Instance, trx *trmgorm.CtxGetter) *Repository {
	if db == nil {
		panic("database instance is nil")
	}

	if trx == nil {
		panic("transaction CtxGetter is nil")
	}

	r := Repository{}

	r.SetTransactionDB(db, trx)

	return &r
}

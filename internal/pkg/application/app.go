package application

import (
	trmgorm "github.com/avito-tech/go-transaction-manager/gorm"
	"github.com/avito-tech/go-transaction-manager/trm"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/db"
)

type App struct {
	DB        *db.Instance
	TrxGetter *trmgorm.CtxGetter
	TxManager trm.Manager
}

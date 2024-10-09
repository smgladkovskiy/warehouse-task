package tx_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	trmgorm "github.com/avito-tech/go-transaction-manager/gorm"
	trmcontext "github.com/avito-tech/go-transaction-manager/trm/context"
	"github.com/avito-tech/go-transaction-manager/trm/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/db"
	trx "github.com/smgladkovskiy/warehouse-task/internal/pkg/tx"
	queryOptions "github.com/smgladkovskiy/warehouse-task/internal/service/entities/query_options"
)

func TestWithTransactionDB_SetTransactionDB(t *testing.T) {
	t.Parallel()

	repo := trx.WithTransactionDB{}
	mockDB, _, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{NowFunc: func() time.Time { return time.Time{} }})
	require.NoError(t, err)

	dbi := db.Instance{Gorm: gormDB}
	trg := trmgorm.NewCtxGetter(trmcontext.DefaultManager)

	assert.Panics(t, func() {
		repo.SetTransactionDB(nil, trg)
	})
	assert.Panics(t, func() {
		repo.SetTransactionDB(&dbi, nil)
	})
}

func TestWithTransactionDB_GetQueryDB(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string
		in   queryOptions.QueryOptionable
		exp  func(ctx context.Context, t *testing.T, in queryOptions.QueryOptionable, dbi *db.Instance, trg *trmgorm.CtxGetter) *gorm.DB
	}

	tcs := []testCase{
		{
			name: "default async db",
			in:   queryOptions.NewBasicQueryOptions(),
			exp: func(ctx context.Context, t *testing.T, in queryOptions.QueryOptionable, dbi *db.Instance, trg *trmgorm.CtxGetter) *gorm.DB {
				t.Helper()

				return dbi.AsyncDB().WithContext(ctx)
			},
		},
		{
			name: "nil qos default async db",
			in:   nil,
			exp: func(ctx context.Context, t *testing.T, in queryOptions.QueryOptionable, dbi *db.Instance, trg *trmgorm.CtxGetter) *gorm.DB {
				t.Helper()

				return dbi.AsyncDB().WithContext(ctx)
			},
		},
		{
			name: "(*queryOptions.BasicQueryOptions)(nil) qos default async db",
			in:   (*queryOptions.BasicQueryOptions)(nil),
			exp: func(ctx context.Context, t *testing.T, in queryOptions.QueryOptionable, dbi *db.Instance, trg *trmgorm.CtxGetter) *gorm.DB {
				t.Helper()

				return dbi.AsyncDB().WithContext(ctx)
			},
		},
		{
			name: "write db for update",
			in:   queryOptions.NewBasicQueryOptions(queryOptions.WithForUpdate[*queryOptions.BasicQueryOptions]()),
			exp: func(ctx context.Context, t *testing.T, in queryOptions.QueryOptionable, dbi *db.Instance, trg *trmgorm.CtxGetter) *gorm.DB {
				t.Helper()

				return trg.TrOrDB(ctx, settings.DefaultCtxKey, dbi.WriteDB()).Clauses(clause.Locking{Strength: "UPDATE"})
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ass := assert.New(t)
			ctx := context.TODO()
			repo := trx.WithTransactionDB{}
			mockDB, _, _ := sqlmock.New()
			dialector := postgres.New(postgres.Config{
				Conn:       mockDB,
				DriverName: "postgres",
			})

			gormDB, err := gorm.Open(dialector, &gorm.Config{NowFunc: func() time.Time { return time.Time{} }})
			require.NoError(t, err)

			dbi := db.Instance{Gorm: gormDB}
			trg := trmgorm.NewCtxGetter(trmcontext.DefaultManager)
			repo.SetTransactionDB(&dbi, trg)

			exp := tc.exp(ctx, t, tc.in, &dbi, trg)

			out := repo.GetQueryDB(ctx, tc.in)

			ass.Equal(exp.Statement.Clauses, out.Statement.Clauses)
		})
	}
}

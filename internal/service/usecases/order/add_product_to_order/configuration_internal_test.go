package addproducttoorder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/checker"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/log"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/now"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/uuid"
	upsertOrder "github.com/smgladkovskiy/warehouse-task/internal/service/commands/order/upsert"
	upsertOrderProduct "github.com/smgladkovskiy/warehouse-task/internal/service/commands/order_product/upsert"
	getOrderByID "github.com/smgladkovskiy/warehouse-task/internal/service/queries/order/get_order"
	getStocks "github.com/smgladkovskiy/warehouse-task/internal/service/queries/order/get_stocks"
	getProduct "github.com/smgladkovskiy/warehouse-task/internal/service/queries/product/get_product"
	usecase "github.com/smgladkovskiy/warehouse-task/internal/service/usecases"
)

func TestConfiguration(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	loggerMock := log.NewLogMock(ctrl)
	nowFunc := now.NewMock(ctrl)
	uuidFunc := uuid.NewMock(ctrl)
	getOrderMock := getOrderByID.NewGetOrderMock(ctrl)
	getProductMock := getProduct.NewGetProductMock(ctrl)
	getStocksMock := getStocks.NewGetStocksMock(ctrl)
	upsertOrderMock := upsertOrder.NewUpsertOrderMock(ctrl)
	upsertOrderProductMock := upsertOrderProduct.NewUpsertOrderProductMock(ctrl)

	cfgs := []usecase.Configuration[*UseCase]{
		usecase.WithLogger[*UseCase](loggerMock),
		usecase.WithNowFunc[*UseCase](nowFunc),
		usecase.WithUUIDFunc[*UseCase](uuidFunc),
		WithGetOrderQuery(getOrderByID.NewQueryHandler(getOrderMock)),
		WithGetProductQuery(getProduct.NewQueryHandler(getProductMock)),
		WithGetStocksQuery(getStocks.NewQueryHandler(getStocksMock)),
		WithUpsertOrderCommand(upsertOrder.NewCommandHandler(upsertOrderMock)),
		WithUpsertOrderProductCommand(upsertOrderProduct.NewCommandHandler(upsertOrderProductMock)),
	}

	f := WithGetOrderQuery(nil)
	uc, err := NewUseCase(f)
	require.Error(t, err)
	assert.Empty(t, uc)

	f = WithGetProductQuery(nil)
	uc, err = NewUseCase(f)
	require.Error(t, err)
	assert.Empty(t, uc)

	f = WithGetStocksQuery(nil)
	uc, err = NewUseCase(f)
	require.Error(t, err)
	assert.Empty(t, uc)

	f = WithUpsertOrderCommand(nil)
	uc, err = NewUseCase(f)
	require.Error(t, err)
	assert.Empty(t, uc)

	f = WithUpsertOrderProductCommand(nil)
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

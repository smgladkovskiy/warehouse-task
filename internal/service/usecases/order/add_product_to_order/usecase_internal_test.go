package addproducttoorder

import (
	"context"
	"testing"
	"time"

	baseUUID "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/checker"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/log"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/now"
	trx "github.com/smgladkovskiy/warehouse-task/internal/pkg/tx"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/uuid"
	upsertOrder "github.com/smgladkovskiy/warehouse-task/internal/service/commands/order/upsert"
	upsertOrderProduct "github.com/smgladkovskiy/warehouse-task/internal/service/commands/order_product/upsert"
	"github.com/smgladkovskiy/warehouse-task/internal/service/entities"
	queryoptions "github.com/smgladkovskiy/warehouse-task/internal/service/entities/query_options"
	vObject "github.com/smgladkovskiy/warehouse-task/internal/service/entities/value_objects"
	getOrderByID "github.com/smgladkovskiy/warehouse-task/internal/service/queries/order/get_order"
	getStocks "github.com/smgladkovskiy/warehouse-task/internal/service/queries/order/get_stocks"
	getProduct "github.com/smgladkovskiy/warehouse-task/internal/service/queries/product/get_product"
	usecase "github.com/smgladkovskiy/warehouse-task/internal/service/usecases"
)

func TestNewUseCase(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	loggerMock := log.NewLogMock(ctrl)
	txManagerMock := trx.NewTransactionManagerMock(ctrl)
	nowFunc := now.NewMock(ctrl)
	uuidFunc := uuid.NewMock(ctrl)
	getOrderMock := getOrderByID.NewGetOrderMock(ctrl)
	getProductMock := getProduct.NewGetProductMock(ctrl)
	getStocksMock := getStocks.NewGetStocksMock(ctrl)
	upsertOrderMock := upsertOrder.NewUpsertOrderMock(ctrl)
	upsertOrderProductMock := upsertOrderProduct.NewUpsertOrderProductMock(ctrl)

	cfgs := []usecase.Configuration[*UseCase]{
		usecase.WithTransactionManager[*UseCase](txManagerMock),
		usecase.WithLogger[*UseCase](loggerMock),
		usecase.WithNowFunc[*UseCase](nowFunc),
		usecase.WithUUIDFunc[*UseCase](uuidFunc),
		WithGetOrderQuery(getOrderByID.NewQueryHandler(getOrderMock)),
		WithGetProductQuery(getProduct.NewQueryHandler(getProductMock)),
		WithGetStocksQuery(getStocks.NewQueryHandler(getStocksMock)),
		WithUpsertOrderCommand(upsertOrder.NewCommandHandler(upsertOrderMock)),
		WithUpsertOrderProductCommand(upsertOrderProduct.NewCommandHandler(upsertOrderProductMock)),
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
		exp  func(t *testing.T, in testRequest, loggerMock *log.LogMock, trxMng *trx.TransactionManagerMock) error
	}

	tcs := []testCase{
		{
			name: "happy path",
			in:   testRequest{},
			exp: func(t *testing.T, in testRequest, loggerMock *log.LogMock, trxMng *trx.TransactionManagerMock) error {
				trxMng.EXPECT().Do(gomock.Any(), gomock.Any()).Return(nil)
				loggerMock.EXPECT().Debug(gomock.Any(), "END usecase")

				return nil
			},
		},
		{
			name: "transaction error",
			in:   testRequest{},
			exp: func(t *testing.T, in testRequest, loggerMock *log.LogMock, trxMng *trx.TransactionManagerMock) error {
				trxMng.EXPECT().Do(gomock.Any(), gomock.Any()).Return(assert.AnError)

				loggerMock.EXPECT().Error(gomock.Any(), "STOP usecase! transaction error", log.Err(assert.AnError))

				return assert.AnError
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			loggerMock := log.NewLogMock(ctrl)
			txManagerMock := trx.NewTransactionManagerMock(ctrl)
			nowFunc := now.NewMock(ctrl)
			uuidFunc := uuid.NewMock(ctrl)
			getOrderMock := getOrderByID.NewGetOrderMock(ctrl)
			getProductMock := getProduct.NewGetProductMock(ctrl)
			getStocksMock := getStocks.NewGetStocksMock(ctrl)
			upsertOrderMock := upsertOrder.NewUpsertOrderMock(ctrl)
			upsertOrderProductMock := upsertOrderProduct.NewUpsertOrderProductMock(ctrl)

			cfgs := []usecase.Configuration[*UseCase]{
				usecase.WithTransactionManager[*UseCase](txManagerMock),
				usecase.WithLogger[*UseCase](loggerMock),
				usecase.WithNowFunc[*UseCase](nowFunc),
				usecase.WithUUIDFunc[*UseCase](uuidFunc),
				WithGetOrderQuery(getOrderByID.NewQueryHandler(getOrderMock)),
				WithGetProductQuery(getProduct.NewQueryHandler(getProductMock)),
				WithGetStocksQuery(getStocks.NewQueryHandler(getStocksMock)),
				WithUpsertOrderCommand(upsertOrder.NewCommandHandler(upsertOrderMock)),
				WithUpsertOrderProductCommand(upsertOrderProduct.NewCommandHandler(upsertOrderProductMock)),
			}

			loggerMock.EXPECT().With(
				log.String("orderUUID", tc.in.GetOrderID().String()),
				log.String("userUUID", tc.in.GetUserID().String()),
				log.String("productUUID", tc.in.GetProductID().String()),
				log.Uint64("quantity", tc.in.GetQuantity()),
			).Return(loggerMock)
			loggerMock.EXPECT().Debug(gomock.Any(), "START usecase")

			uc, err := NewUseCase(cfgs...)
			require.NoError(t, err)

			expErr := tc.exp(t, tc.in, loggerMock, txManagerMock)

			assert.ErrorIs(t, uc.Run(context.Background(), tc.in), expErr)
		})
	}
}

func TestUseCase_transaction(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string
		in   testRequest
		exp  func(t *testing.T, in testRequest, loggerMock *log.LogMock, getOrderMock *getOrderByID.GetOrderMock, getProductMock *getProduct.GetProductMock, getStocksMock *getStocks.GetStocksMock, upsertOrderMock *upsertOrder.UpsertOrderMock, upsertOrderProductMock *upsertOrderProduct.UpsertOrderProductMock) error
	}

	tn := time.Now()
	id := baseUUID.New()

	nowFunc := now.NewMock(gomock.NewController(t))
	uuidFunc := uuid.NewMock(gomock.NewController(t))

	nowFunc.EXPECT().Now().AnyTimes().Return(tn)
	nowFunc.EXPECT().NowP().AnyTimes().Return(&tn)
	uuidFunc.EXPECT().UUID().AnyTimes().Return(id)

	tcs := []testCase{
		{
			name: "happy path",
			in: testRequest{
				orderUUID:   id,
				productUUID: id,
				quantity:    6,
				userUUID:    id,
			},
			exp: func(t *testing.T, in testRequest, loggerMock *log.LogMock, getOrderMock *getOrderByID.GetOrderMock, getProductMock *getProduct.GetProductMock, getStocksMock *getStocks.GetStocksMock, upsertOrderMock *upsertOrder.UpsertOrderMock, upsertOrderProductMock *upsertOrderProduct.UpsertOrderProductMock) error {
				t.Helper()

				order := entities.NewOrderUnsafe(
					vObject.NewUserIDFromUUIDUnsafe(in.GetOrderID()),
					entities.WithUUIDFunc[*entities.Order](uuidFunc),
					entities.WithNowFunc[*entities.Order](nowFunc),
				)
				product := entities.NewProductUnsafe(
					vObject.NewProductTitleUnsafe("product title"),
					vObject.NewProductDescriptionUnsafe("product description"),
					vObject.NewPriceUnsafe(10000),
					entities.WithUUIDFunc[*entities.Product](uuidFunc),
					entities.WithNowFunc[*entities.Product](nowFunc),
				)
				productStocks := entities.Stocks{
					entities.NewStockUnsafe(
						product.ID,
						vObject.NewWarehouseIDFromUUIDUnsafe(baseUUID.New()),
						vObject.NewQuantityUnsafe(3),  //reserve
						vObject.NewQuantityUnsafe(10), //available
					),
					entities.NewStockUnsafe(
						product.ID,
						vObject.NewWarehouseIDFromUUIDUnsafe(baseUUID.New()),
						vObject.NewQuantityUnsafe(5),   //reserve
						vObject.NewQuantityUnsafe(100), //available
					),
				}

				getOrderMock.EXPECT().GetOrder(gomock.Any(), queryoptions.NewOrderQueryOptions(queryoptions.WithOrderID(order.ID), queryoptions.WithForUpdate[*queryoptions.OrderQueryOptions]())).Return(&order, nil)
				loggerMock.EXPECT().With(log.String("orderID", order.ID.String())).Return(loggerMock)
				getProductMock.EXPECT().GetProduct(gomock.Any(), queryoptions.NewProductQueryOptions(queryoptions.WithProductID(product.ID))).Return(&product, nil)
				getStocksMock.EXPECT().GetStocks(gomock.Any(), queryoptions.NewStockQueryOptions(queryoptions.WithStockProductID(product.ID))).Return(productStocks, nil)
				loggerMock.EXPECT().With(log.Uint64("productAvailableQuantity", productStocks.GetAvailableQuantity().Uint64())).Return(loggerMock)

				changedOrder := order
				require.NoError(t, changedOrder.ChangeOrderProducts(productStocks, product, in.GetQuantity()))

				upsertOrderMock.EXPECT().UpsertOrder(gomock.Any(), &changedOrder).Return(nil)

				orderProduct := changedOrder.GetOrderProductByProductIDUnsafe(product.ID)

				upsertOrderProductMock.EXPECT().UpsertOrderProduct(gomock.Any(), orderProduct).Return(nil)
				loggerMock.EXPECT().With(log.Uint64("orderProductQuantity", orderProduct.Quantity.Uint64())).Return(loggerMock)

				return nil
			},
		},
		{
			name: "orderProduct upsert error",
			in: testRequest{
				orderUUID:   id,
				productUUID: id,
				quantity:    6,
				userUUID:    id,
			},
			exp: func(t *testing.T, in testRequest, loggerMock *log.LogMock, getOrderMock *getOrderByID.GetOrderMock, getProductMock *getProduct.GetProductMock, getStocksMock *getStocks.GetStocksMock, upsertOrderMock *upsertOrder.UpsertOrderMock, upsertOrderProductMock *upsertOrderProduct.UpsertOrderProductMock) error {
				t.Helper()

				order := entities.NewOrderUnsafe(
					vObject.NewUserIDFromUUIDUnsafe(in.GetOrderID()),
					entities.WithUUIDFunc[*entities.Order](uuidFunc),
					entities.WithNowFunc[*entities.Order](nowFunc),
				)
				product := entities.NewProductUnsafe(
					vObject.NewProductTitleUnsafe("product title"),
					vObject.NewProductDescriptionUnsafe("product description"),
					vObject.NewPriceUnsafe(10000),
					entities.WithUUIDFunc[*entities.Product](uuidFunc),
					entities.WithNowFunc[*entities.Product](nowFunc),
				)
				productStocks := entities.Stocks{
					entities.NewStockUnsafe(
						product.ID,
						vObject.NewWarehouseIDFromUUIDUnsafe(baseUUID.New()),
						vObject.NewQuantityUnsafe(3),  //reserve
						vObject.NewQuantityUnsafe(10), //available
					),
					entities.NewStockUnsafe(
						product.ID,
						vObject.NewWarehouseIDFromUUIDUnsafe(baseUUID.New()),
						vObject.NewQuantityUnsafe(5),   //reserve
						vObject.NewQuantityUnsafe(100), //available
					),
				}

				getOrderMock.EXPECT().GetOrder(gomock.Any(), queryoptions.NewOrderQueryOptions(queryoptions.WithOrderID(order.ID), queryoptions.WithForUpdate[*queryoptions.OrderQueryOptions]())).Return(&order, nil)
				loggerMock.EXPECT().With(log.String("orderID", order.ID.String())).Return(loggerMock)
				getProductMock.EXPECT().GetProduct(gomock.Any(), queryoptions.NewProductQueryOptions(queryoptions.WithProductID(product.ID))).Return(&product, nil)
				getStocksMock.EXPECT().GetStocks(gomock.Any(), queryoptions.NewStockQueryOptions(queryoptions.WithStockProductID(product.ID))).Return(productStocks, nil)
				loggerMock.EXPECT().With(log.Uint64("productAvailableQuantity", productStocks.GetAvailableQuantity().Uint64())).Return(loggerMock)

				changedOrder := order
				require.NoError(t, changedOrder.ChangeOrderProducts(productStocks, product, in.GetQuantity()))

				upsertOrderMock.EXPECT().UpsertOrder(gomock.Any(), &changedOrder).Return(nil)

				orderProduct := changedOrder.GetOrderProductByProductIDUnsafe(product.ID)

				upsertOrderProductMock.EXPECT().UpsertOrderProduct(gomock.Any(), orderProduct).Return(assert.AnError)

				return assert.AnError
			},
		},
		{
			name: "order upsert error",
			in: testRequest{
				orderUUID:   id,
				productUUID: id,
				quantity:    6,
				userUUID:    id,
			},
			exp: func(t *testing.T, in testRequest, loggerMock *log.LogMock, getOrderMock *getOrderByID.GetOrderMock, getProductMock *getProduct.GetProductMock, getStocksMock *getStocks.GetStocksMock, upsertOrderMock *upsertOrder.UpsertOrderMock, upsertOrderProductMock *upsertOrderProduct.UpsertOrderProductMock) error {
				t.Helper()

				order := entities.NewOrderUnsafe(
					vObject.NewUserIDFromUUIDUnsafe(in.GetOrderID()),
					entities.WithUUIDFunc[*entities.Order](uuidFunc),
					entities.WithNowFunc[*entities.Order](nowFunc),
				)
				product := entities.NewProductUnsafe(
					vObject.NewProductTitleUnsafe("product title"),
					vObject.NewProductDescriptionUnsafe("product description"),
					vObject.NewPriceUnsafe(10000),
					entities.WithUUIDFunc[*entities.Product](uuidFunc),
					entities.WithNowFunc[*entities.Product](nowFunc),
				)
				productStocks := entities.Stocks{
					entities.NewStockUnsafe(
						product.ID,
						vObject.NewWarehouseIDFromUUIDUnsafe(baseUUID.New()),
						vObject.NewQuantityUnsafe(3),  //reserve
						vObject.NewQuantityUnsafe(10), //available
					),
					entities.NewStockUnsafe(
						product.ID,
						vObject.NewWarehouseIDFromUUIDUnsafe(baseUUID.New()),
						vObject.NewQuantityUnsafe(5),   //reserve
						vObject.NewQuantityUnsafe(100), //available
					),
				}

				getOrderMock.EXPECT().GetOrder(gomock.Any(), queryoptions.NewOrderQueryOptions(queryoptions.WithOrderID(order.ID), queryoptions.WithForUpdate[*queryoptions.OrderQueryOptions]())).Return(&order, nil)
				loggerMock.EXPECT().With(log.String("orderID", order.ID.String())).Return(loggerMock)
				getProductMock.EXPECT().GetProduct(gomock.Any(), queryoptions.NewProductQueryOptions(queryoptions.WithProductID(product.ID))).Return(&product, nil)
				getStocksMock.EXPECT().GetStocks(gomock.Any(), queryoptions.NewStockQueryOptions(queryoptions.WithStockProductID(product.ID))).Return(productStocks, nil)
				loggerMock.EXPECT().With(log.Uint64("productAvailableQuantity", productStocks.GetAvailableQuantity().Uint64())).Return(loggerMock)

				changedOrder := order
				require.NoError(t, changedOrder.ChangeOrderProducts(productStocks, product, in.GetQuantity()))

				upsertOrderMock.EXPECT().UpsertOrder(gomock.Any(), &changedOrder).Return(assert.AnError)

				return assert.AnError
			},
		},
		{
			name: "stocks are les than requested quantity",
			in: testRequest{
				orderUUID:   id,
				productUUID: id,
				quantity:    111,
				userUUID:    id,
			},
			exp: func(t *testing.T, in testRequest, loggerMock *log.LogMock, getOrderMock *getOrderByID.GetOrderMock, getProductMock *getProduct.GetProductMock, getStocksMock *getStocks.GetStocksMock, upsertOrderMock *upsertOrder.UpsertOrderMock, upsertOrderProductMock *upsertOrderProduct.UpsertOrderProductMock) error {
				t.Helper()

				order := entities.NewOrderUnsafe(
					vObject.NewUserIDFromUUIDUnsafe(in.GetOrderID()),
					entities.WithUUIDFunc[*entities.Order](uuidFunc),
					entities.WithNowFunc[*entities.Order](nowFunc),
				)
				product := entities.NewProductUnsafe(
					vObject.NewProductTitleUnsafe("product title"),
					vObject.NewProductDescriptionUnsafe("product description"),
					vObject.NewPriceUnsafe(10000),
					entities.WithUUIDFunc[*entities.Product](uuidFunc),
					entities.WithNowFunc[*entities.Product](nowFunc),
				)
				productStocks := entities.Stocks{
					entities.NewStockUnsafe(
						product.ID,
						vObject.NewWarehouseIDFromUUIDUnsafe(baseUUID.New()),
						vObject.NewQuantityUnsafe(3),  //reserve
						vObject.NewQuantityUnsafe(10), //available
					),
					entities.NewStockUnsafe(
						product.ID,
						vObject.NewWarehouseIDFromUUIDUnsafe(baseUUID.New()),
						vObject.NewQuantityUnsafe(5),   //reserve
						vObject.NewQuantityUnsafe(100), //available
					),
				}

				getOrderMock.EXPECT().GetOrder(gomock.Any(), queryoptions.NewOrderQueryOptions(queryoptions.WithOrderID(order.ID), queryoptions.WithForUpdate[*queryoptions.OrderQueryOptions]())).Return(&order, nil)
				loggerMock.EXPECT().With(log.String("orderID", order.ID.String())).Return(loggerMock)
				getProductMock.EXPECT().GetProduct(gomock.Any(), queryoptions.NewProductQueryOptions(queryoptions.WithProductID(product.ID))).Return(&product, nil)
				getStocksMock.EXPECT().GetStocks(gomock.Any(), queryoptions.NewStockQueryOptions(queryoptions.WithStockProductID(product.ID))).Return(productStocks, nil)
				loggerMock.EXPECT().With(log.Uint64("productAvailableQuantity", productStocks.GetAvailableQuantity().Uint64())).Return(loggerMock)

				return entities.ErrNotEnoughProductIntStocks
			},
		},
		{
			name: "get product stocks error",
			in: testRequest{
				orderUUID:   id,
				productUUID: id,
				quantity:    111,
				userUUID:    id,
			},
			exp: func(t *testing.T, in testRequest, loggerMock *log.LogMock, getOrderMock *getOrderByID.GetOrderMock, getProductMock *getProduct.GetProductMock, getStocksMock *getStocks.GetStocksMock, upsertOrderMock *upsertOrder.UpsertOrderMock, upsertOrderProductMock *upsertOrderProduct.UpsertOrderProductMock) error {
				t.Helper()

				order := entities.NewOrderUnsafe(
					vObject.NewUserIDFromUUIDUnsafe(in.GetOrderID()),
					entities.WithUUIDFunc[*entities.Order](uuidFunc),
					entities.WithNowFunc[*entities.Order](nowFunc),
				)
				product := entities.NewProductUnsafe(
					vObject.NewProductTitleUnsafe("product title"),
					vObject.NewProductDescriptionUnsafe("product description"),
					vObject.NewPriceUnsafe(10000),
					entities.WithUUIDFunc[*entities.Product](uuidFunc),
					entities.WithNowFunc[*entities.Product](nowFunc),
				)

				getOrderMock.EXPECT().GetOrder(gomock.Any(), queryoptions.NewOrderQueryOptions(queryoptions.WithOrderID(order.ID), queryoptions.WithForUpdate[*queryoptions.OrderQueryOptions]())).Return(&order, nil)
				loggerMock.EXPECT().With(log.String("orderID", order.ID.String())).Return(loggerMock)
				getProductMock.EXPECT().GetProduct(gomock.Any(), queryoptions.NewProductQueryOptions(queryoptions.WithProductID(product.ID))).Return(&product, nil)
				getStocksMock.EXPECT().GetStocks(gomock.Any(), queryoptions.NewStockQueryOptions(queryoptions.WithStockProductID(product.ID))).Return(nil, assert.AnError)

				return assert.AnError
			},
		},
		{
			name: "get product handler error",
			in: testRequest{
				orderUUID:   id,
				productUUID: id,
				quantity:    111,
				userUUID:    id,
			},
			exp: func(t *testing.T, in testRequest, loggerMock *log.LogMock, getOrderMock *getOrderByID.GetOrderMock, getProductMock *getProduct.GetProductMock, getStocksMock *getStocks.GetStocksMock, upsertOrderMock *upsertOrder.UpsertOrderMock, upsertOrderProductMock *upsertOrderProduct.UpsertOrderProductMock) error {
				t.Helper()

				order := entities.NewOrderUnsafe(
					vObject.NewUserIDFromUUIDUnsafe(in.GetOrderID()),
					entities.WithUUIDFunc[*entities.Order](uuidFunc),
					entities.WithNowFunc[*entities.Order](nowFunc),
				)
				product := entities.NewProductUnsafe(
					vObject.NewProductTitleUnsafe("product title"),
					vObject.NewProductDescriptionUnsafe("product description"),
					vObject.NewPriceUnsafe(10000),
					entities.WithUUIDFunc[*entities.Product](uuidFunc),
					entities.WithNowFunc[*entities.Product](nowFunc),
				)

				getOrderMock.EXPECT().GetOrder(gomock.Any(), queryoptions.NewOrderQueryOptions(queryoptions.WithOrderID(order.ID), queryoptions.WithForUpdate[*queryoptions.OrderQueryOptions]())).Return(&order, nil)
				loggerMock.EXPECT().With(log.String("orderID", order.ID.String())).Return(loggerMock)
				getProductMock.EXPECT().GetProduct(gomock.Any(), queryoptions.NewProductQueryOptions(queryoptions.WithProductID(product.ID))).Return(nil, assert.AnError)

				return assert.AnError
			},
		},
		{
			name: "get product query error",
			in: testRequest{
				orderUUID:   id,
				productUUID: baseUUID.Nil,
				quantity:    111,
				userUUID:    id,
			},
			exp: func(t *testing.T, in testRequest, loggerMock *log.LogMock, getOrderMock *getOrderByID.GetOrderMock, getProductMock *getProduct.GetProductMock, getStocksMock *getStocks.GetStocksMock, upsertOrderMock *upsertOrder.UpsertOrderMock, upsertOrderProductMock *upsertOrderProduct.UpsertOrderProductMock) error {
				t.Helper()

				order := entities.NewOrderUnsafe(
					vObject.NewUserIDFromUUIDUnsafe(in.GetOrderID()),
					entities.WithUUIDFunc[*entities.Order](uuidFunc),
					entities.WithNowFunc[*entities.Order](nowFunc),
				)

				getOrderMock.EXPECT().GetOrder(gomock.Any(), queryoptions.NewOrderQueryOptions(queryoptions.WithOrderID(order.ID), queryoptions.WithForUpdate[*queryoptions.OrderQueryOptions]())).Return(&order, nil)
				loggerMock.EXPECT().With(log.String("orderID", order.ID.String())).Return(loggerMock)

				return vObject.ErrEmptyID
			},
		},
		{
			name: "get order error",
			in: testRequest{
				orderUUID:   id,
				productUUID: id,
				quantity:    111,
				userUUID:    id,
			},
			exp: func(t *testing.T, in testRequest, loggerMock *log.LogMock, getOrderMock *getOrderByID.GetOrderMock, getProductMock *getProduct.GetProductMock, getStocksMock *getStocks.GetStocksMock, upsertOrderMock *upsertOrder.UpsertOrderMock, upsertOrderProductMock *upsertOrderProduct.UpsertOrderProductMock) error {
				t.Helper()

				order := entities.NewOrderUnsafe(
					vObject.NewUserIDFromUUIDUnsafe(in.GetOrderID()),
					entities.WithUUIDFunc[*entities.Order](uuidFunc),
					entities.WithNowFunc[*entities.Order](nowFunc),
				)

				getOrderMock.EXPECT().GetOrder(gomock.Any(), queryoptions.NewOrderQueryOptions(queryoptions.WithOrderID(order.ID), queryoptions.WithForUpdate[*queryoptions.OrderQueryOptions]())).Return(nil, assert.AnError)

				return assert.AnError
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			loggerMock := log.NewLogMock(ctrl)
			txManagerMock := trx.NewTransactionManagerMock(ctrl)
			getOrderMock := getOrderByID.NewGetOrderMock(ctrl)
			getProductMock := getProduct.NewGetProductMock(ctrl)
			getStocksMock := getStocks.NewGetStocksMock(ctrl)
			upsertOrderMock := upsertOrder.NewUpsertOrderMock(ctrl)
			upsertOrderProductMock := upsertOrderProduct.NewUpsertOrderProductMock(ctrl)

			cfgs := []usecase.Configuration[*UseCase]{
				usecase.WithTransactionManager[*UseCase](txManagerMock),
				usecase.WithLogger[*UseCase](loggerMock),
				usecase.WithNowFunc[*UseCase](nowFunc),
				usecase.WithUUIDFunc[*UseCase](uuidFunc),
				WithGetOrderQuery(getOrderByID.NewQueryHandler(getOrderMock)),
				WithGetProductQuery(getProduct.NewQueryHandler(getProductMock)),
				WithGetStocksQuery(getStocks.NewQueryHandler(getStocksMock)),
				WithUpsertOrderCommand(upsertOrder.NewCommandHandler(upsertOrderMock)),
				WithUpsertOrderProductCommand(upsertOrderProduct.NewCommandHandler(upsertOrderProductMock)),
			}

			uc, err := NewUseCase(cfgs...)
			require.NoError(t, err)

			expErr := tc.exp(t, tc.in, loggerMock, getOrderMock, getProductMock, getStocksMock, upsertOrderMock, upsertOrderProductMock)

			errFunc := uc.transaction(loggerMock, tc.in)

			assert.ErrorIs(t, errFunc(context.Background()), expErr)
		})
	}
}

func TestUseCase_getOrder(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string
		in   testRequest
		exp  func(t *testing.T, in testRequest, getOrderMock *getOrderByID.GetOrderMock, upsertOrderMock *upsertOrder.UpsertOrderMock) (*entities.Order, error)
	}

	tn := time.Now()
	id := baseUUID.New()

	nowFunc := now.NewMock(gomock.NewController(t))
	uuidFunc := uuid.NewMock(gomock.NewController(t))

	nowFunc.EXPECT().Now().AnyTimes().Return(tn)
	nowFunc.EXPECT().NowP().AnyTimes().Return(&tn)
	uuidFunc.EXPECT().UUID().AnyTimes().Return(id)

	tcs := []testCase{
		{
			name: "get order happy path",
			in: testRequest{
				orderUUID: id,
			},
			exp: func(t *testing.T, in testRequest, getOrderMock *getOrderByID.GetOrderMock, upsertOrderMock *upsertOrder.UpsertOrderMock) (*entities.Order, error) {
				t.Helper()

				order := entities.NewOrderUnsafe(
					vObject.NewUserIDFromUUIDUnsafe(in.GetOrderID()),
					entities.WithUUIDFunc[*entities.Order](uuidFunc),
					entities.WithNowFunc[*entities.Order](nowFunc),
				)

				getOrderMock.EXPECT().GetOrder(gomock.Any(), queryoptions.NewOrderQueryOptions(queryoptions.WithOrderID(order.ID), queryoptions.WithForUpdate[*queryoptions.OrderQueryOptions]())).Return(&order, nil)

				return &order, nil
			},
		},
		{
			name: "create order happy path",
			in: testRequest{
				orderUUID: baseUUID.Nil,
				userUUID:  baseUUID.New(),
			},
			exp: func(t *testing.T, in testRequest, getOrderMock *getOrderByID.GetOrderMock, upsertOrderMock *upsertOrder.UpsertOrderMock) (*entities.Order, error) {
				t.Helper()

				order := entities.NewOrderUnsafe(
					vObject.NewUserIDFromUUIDUnsafe(in.GetUserID()),
					entities.WithUUIDFunc[*entities.Order](uuidFunc),
					entities.WithNowFunc[*entities.Order](nowFunc),
				)

				upsertOrderMock.EXPECT().UpsertOrder(gomock.Any(), &order).Return(nil)

				return &order, nil
			},
		},
		{
			name: "get order handler error",
			in: testRequest{
				orderUUID: id,
				userUUID:  baseUUID.New(),
			},
			exp: func(t *testing.T, in testRequest, getOrderMock *getOrderByID.GetOrderMock, upsertOrderMock *upsertOrder.UpsertOrderMock) (*entities.Order, error) {
				t.Helper()

				order := entities.NewOrderUnsafe(
					vObject.NewUserIDFromUUIDUnsafe(in.GetUserID()),
					entities.WithUUIDFunc[*entities.Order](uuidFunc),
					entities.WithNowFunc[*entities.Order](nowFunc),
				)

				getOrderMock.EXPECT().GetOrder(gomock.Any(), queryoptions.NewOrderQueryOptions(queryoptions.WithOrderID(order.ID), queryoptions.WithForUpdate[*queryoptions.OrderQueryOptions]())).Return(nil, assert.AnError)

				return nil, assert.AnError
			},
		},
		{
			name: "no order found error",
			in: testRequest{
				orderUUID: id,
				userUUID:  baseUUID.New(),
			},
			exp: func(t *testing.T, in testRequest, getOrderMock *getOrderByID.GetOrderMock, upsertOrderMock *upsertOrder.UpsertOrderMock) (*entities.Order, error) {
				t.Helper()

				order := entities.NewOrderUnsafe(
					vObject.NewUserIDFromUUIDUnsafe(in.GetUserID()),
					entities.WithUUIDFunc[*entities.Order](uuidFunc),
					entities.WithNowFunc[*entities.Order](nowFunc),
				)

				getOrderMock.EXPECT().GetOrder(gomock.Any(), queryoptions.NewOrderQueryOptions(queryoptions.WithOrderID(order.ID), queryoptions.WithForUpdate[*queryoptions.OrderQueryOptions]())).Return(nil, entities.ErrOrderRecNotFound)
				upsertOrderMock.EXPECT().UpsertOrder(gomock.Any(), &order).Return(nil)

				return &order, nil
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			loggerMock := log.NewLogMock(ctrl)
			txManagerMock := trx.NewTransactionManagerMock(ctrl)
			getOrderMock := getOrderByID.NewGetOrderMock(ctrl)
			getProductMock := getProduct.NewGetProductMock(ctrl)
			getStocksMock := getStocks.NewGetStocksMock(ctrl)
			upsertOrderMock := upsertOrder.NewUpsertOrderMock(ctrl)
			upsertOrderProductMock := upsertOrderProduct.NewUpsertOrderProductMock(ctrl)

			cfgs := []usecase.Configuration[*UseCase]{
				usecase.WithTransactionManager[*UseCase](txManagerMock),
				usecase.WithLogger[*UseCase](loggerMock),
				usecase.WithNowFunc[*UseCase](nowFunc),
				usecase.WithUUIDFunc[*UseCase](uuidFunc),
				WithGetOrderQuery(getOrderByID.NewQueryHandler(getOrderMock)),
				WithGetProductQuery(getProduct.NewQueryHandler(getProductMock)),
				WithGetStocksQuery(getStocks.NewQueryHandler(getStocksMock)),
				WithUpsertOrderCommand(upsertOrder.NewCommandHandler(upsertOrderMock)),
				WithUpsertOrderProductCommand(upsertOrderProduct.NewCommandHandler(upsertOrderProductMock)),
			}

			uc, err := NewUseCase(cfgs...)
			require.NoError(t, err)

			expOut, expErr := tc.exp(t, tc.in, getOrderMock, upsertOrderMock)

			out, err := uc.getOrder(context.TODO(), tc.in)

			assert.ErrorIs(t, err, expErr)
			assert.Equal(t, out, expOut)
		})
	}
}

func TestUseCase_createOrder(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string
		in   testRequest
		exp  func(t *testing.T, in testRequest, getOrderMock *getOrderByID.GetOrderMock, upsertOrderMock *upsertOrder.UpsertOrderMock) (*entities.Order, error)
	}

	tn := time.Now()
	id := baseUUID.New()

	nowFunc := now.NewMock(gomock.NewController(t))
	uuidFunc := uuid.NewMock(gomock.NewController(t))

	nowFunc.EXPECT().Now().AnyTimes().Return(tn)
	nowFunc.EXPECT().NowP().AnyTimes().Return(&tn)
	uuidFunc.EXPECT().UUID().AnyTimes().Return(id)

	tcs := []testCase{
		{
			name: "create order happy path",
			in: testRequest{
				orderUUID: baseUUID.Nil,
				userUUID:  baseUUID.New(),
			},
			exp: func(t *testing.T, in testRequest, getOrderMock *getOrderByID.GetOrderMock, upsertOrderMock *upsertOrder.UpsertOrderMock) (*entities.Order, error) {
				t.Helper()

				order := entities.NewOrderUnsafe(
					vObject.NewUserIDFromUUIDUnsafe(in.GetUserID()),
					entities.WithUUIDFunc[*entities.Order](uuidFunc),
					entities.WithNowFunc[*entities.Order](nowFunc),
				)

				upsertOrderMock.EXPECT().UpsertOrder(gomock.Any(), &order).Return(nil)

				return &order, nil
			},
		},
		{
			name: "upsert order error",
			in: testRequest{
				orderUUID: baseUUID.Nil,
				userUUID:  baseUUID.New(),
			},
			exp: func(t *testing.T, in testRequest, getOrderMock *getOrderByID.GetOrderMock, upsertOrderMock *upsertOrder.UpsertOrderMock) (*entities.Order, error) {
				t.Helper()

				order := entities.NewOrderUnsafe(
					vObject.NewUserIDFromUUIDUnsafe(in.GetUserID()),
					entities.WithUUIDFunc[*entities.Order](uuidFunc),
					entities.WithNowFunc[*entities.Order](nowFunc),
				)

				upsertOrderMock.EXPECT().UpsertOrder(gomock.Any(), &order).Return(assert.AnError)

				return nil, assert.AnError
			},
		},
		{
			name: "upsert order command error",
			in: testRequest{
				orderUUID: baseUUID.Nil,
				userUUID:  baseUUID.Nil,
			},
			exp: func(t *testing.T, in testRequest, getOrderMock *getOrderByID.GetOrderMock, upsertOrderMock *upsertOrder.UpsertOrderMock) (*entities.Order, error) {
				t.Helper()

				return nil, vObject.ErrEmptyID
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			loggerMock := log.NewLogMock(ctrl)
			txManagerMock := trx.NewTransactionManagerMock(ctrl)
			getOrderMock := getOrderByID.NewGetOrderMock(ctrl)
			getProductMock := getProduct.NewGetProductMock(ctrl)
			getStocksMock := getStocks.NewGetStocksMock(ctrl)
			upsertOrderMock := upsertOrder.NewUpsertOrderMock(ctrl)
			upsertOrderProductMock := upsertOrderProduct.NewUpsertOrderProductMock(ctrl)

			cfgs := []usecase.Configuration[*UseCase]{
				usecase.WithTransactionManager[*UseCase](txManagerMock),
				usecase.WithLogger[*UseCase](loggerMock),
				usecase.WithNowFunc[*UseCase](nowFunc),
				usecase.WithUUIDFunc[*UseCase](uuidFunc),
				WithGetOrderQuery(getOrderByID.NewQueryHandler(getOrderMock)),
				WithGetProductQuery(getProduct.NewQueryHandler(getProductMock)),
				WithGetStocksQuery(getStocks.NewQueryHandler(getStocksMock)),
				WithUpsertOrderCommand(upsertOrder.NewCommandHandler(upsertOrderMock)),
				WithUpsertOrderProductCommand(upsertOrderProduct.NewCommandHandler(upsertOrderProductMock)),
			}

			expOut, expErr := tc.exp(t, tc.in, getOrderMock, upsertOrderMock)

			uc, err := NewUseCase(cfgs...)
			require.NoError(t, err)

			out, err := uc.createOrder(context.TODO(), tc.in)

			assert.ErrorIs(t, err, expErr)
			assert.Equal(t, out, expOut)
		})
	}
}

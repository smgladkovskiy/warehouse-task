package entities

import (
	"errors"
	"fmt"
	"time"

	baseUUID "github.com/google/uuid"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/now"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/uuid"
	vObject "github.com/smgladkovskiy/warehouse-task/internal/service/entities/value_objects"
)

type Order struct {
	now.WithNowGenerator
	uuid.WithUUIDGenerator

	ID         vObject.OrderID
	UserID     vObject.UserID
	Status     vObject.OrderStatus
	TotalPrice vObject.Price
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time

	User     *User
	Products OrderProducts
}

var ErrOrderRecNotFound = errors.New("order record not found")

func NewOrder(userUUID baseUUID.UUID, opts ...Option[*Order]) (*Order, error) {
	userID, err := vObject.NewUserIDFromUUID(userUUID)
	if err != nil {
		return nil, fmt.Errorf("[NewOrder - vObject.NewUserIDFromUUID error]: %w", err)
	}

	o := Order{
		UserID: userID,
		Status: vObject.OrderStatusCreated,
	}

	for _, opt := range opts {
		if err = opt(&o); err != nil {
			return nil, fmt.Errorf("[NewOrder - opt error]: %w", err)
		}
	}

	o.ID = vObject.NewOrderIDFromUUIDUnsafe(o.UUID())
	tn := o.Now()
	o.CreatedAt = tn
	o.UpdatedAt = tn

	return &o, nil
}

func NewOrderUnsafe(userID vObject.UserID, opts ...Option[*Order]) Order {
	o, _ := NewOrder(userID.UUID(), opts...)

	return *o
}

func (o *Order) ChangeOrderProducts(stocks Stocks, product Product, quantity uint64) error {
	if stocks.GetAvailableQuantity().IsLessThan(quantity) {
		return fmt.Errorf("[Order.ChangeOrderProducts error]: %w", ErrNotEnoughProductIntStocks)
	}

	orderProduct := o.GetOrCreateOrderProductByProduct(product)

	if quantity == 0 {
		o.Products.Delete(orderProduct)
		o.TotalPrice.Subtract(orderProduct.TotalPrice())
	} else {
		o.TotalPrice.Subtract(orderProduct.TotalPrice())
		orderProduct.ChangeQuantity(quantity)
		o.TotalPrice.Add(orderProduct.TotalPrice())
		o.Products.Replace(*orderProduct)
	}

	return nil
}

func (o *Order) GetOrderProductByProductIDUnsafe(productID vObject.ProductID) *OrderProduct {
	for _, orderProduct := range o.Products {
		if orderProduct.ProductID == productID {
			return &orderProduct
		}
	}

	return nil
}

func (o *Order) GetOrCreateOrderProductByProduct(product Product) *OrderProduct {
	for _, orderProduct := range o.Products {
		if orderProduct.ProductID == product.ID {
			return &orderProduct
		}
	}

	op := NewOrderProductUnsafe(
		o.ID,
		product.ID,
		WithOrderProductPrice(product.Price),
		WithNowFunc[*OrderProduct](o.GetNowGen()),
	)
	op.Product = &product

	o.Products = append(o.Products, op)

	return &op
}

package entities

import (
	"errors"
	"time"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/now"
	vObject "github.com/smgladkovskiy/warehouse-task/internal/service/entities/value_objects"
)

type Stock struct {
	now.WithNowGenerator

	ProductID         vObject.ProductID
	WarehouseID       vObject.WarehouseID
	AvailableQuantity vObject.Quantity
	ReservedQuantity  vObject.Quantity
	CreatedAt         time.Time
}

type Stocks []Stock

var ErrNotEnoughProductIntStocks = errors.New("not enough products in stocks")

func (s Stocks) GetAvailableQuantity() vObject.Quantity {
	var quantity vObject.Quantity

	for _, stock := range s {
		quantity += stock.AvailableQuantity
		quantity -= stock.ReservedQuantity
	}

	return quantity
}

func (s Stocks) GetProductID() vObject.ProductID {
	return s[0].ProductID
}

func NewStockUnsafe(
	productID vObject.ProductID,
	warehouseID vObject.WarehouseID,
	reservedQuantity, availableQuantity vObject.Quantity,
	opts ...Option[*Stock],
) Stock {
	s := Stock{
		ProductID:         productID,
		WarehouseID:       warehouseID,
		ReservedQuantity:  reservedQuantity,
		AvailableQuantity: availableQuantity,
	}

	for _, opt := range opts {
		_ = opt(&s)
	}

	s.CreatedAt = s.Now()

	return s
}

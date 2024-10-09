package entities

import (
	"time"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/now"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/uuid"
	vObject "github.com/smgladkovskiy/warehouse-task/internal/service/entities/value_objects"
)

type Product struct {
	uuid.WithUUIDGenerator
	now.WithNowGenerator

	ID          vObject.ProductID
	Title       vObject.ProductTitle
	Description vObject.ProductDescription
	Tags        vObject.Tags
	Price       vObject.Price
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time

	Remains   Stocks
	Movements ProductMovements
	Orders    OrderProducts
}

func NewProductUnsafe(title vObject.ProductTitle, description vObject.ProductDescription, price vObject.Price, opts ...Option[*Product]) Product {
	p := Product{
		Title:       title,
		Description: description,
		Price:       price,
	}

	for _, opt := range opts {
		_ = opt(&p)
	}

	p.ID = vObject.NewProductIDFromUUIDUnsafe(p.UUID())
	tn := time.Now()
	p.CreatedAt = tn
	p.UpdatedAt = tn

	return p
}

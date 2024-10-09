package entities

import (
	"time"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/now"
	vObject "github.com/smgladkovskiy/warehouse-task/internal/service/entities/value_objects"
)

type OrderProduct struct {
	now.WithNowGenerator

	OrderID   vObject.OrderID
	ProductID vObject.ProductID
	Quantity  vObject.Quantity
	Price     vObject.Price
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Order   *Order
	Product *Product
}

func (p *OrderProduct) ChangeQuantity(quantity uint64) {
	p.Quantity = vObject.NewQuantityUnsafe(quantity)
	p.UpdatedAt = p.Now()
}

func (p *OrderProduct) Delete() {
	tn := p.Now()
	p.UpdatedAt = tn
	p.DeletedAt = &tn
}

func (p *OrderProduct) TotalPrice() vObject.Price {
	return p.Price.Multiply(p.Quantity)
}

type OrderProducts []OrderProduct

func (p OrderProducts) Delete(product *OrderProduct) {
	for i, orderProduct := range p {
		if orderProduct.ProductID == product.ProductID {
			orderProduct.Delete()
			p[i] = orderProduct
		}
	}
}

func (p OrderProducts) Replace(orderProduct OrderProduct) {
	for i, op := range p {
		if op.ProductID == orderProduct.ProductID {
			p[i] = orderProduct
		}
	}
}

func NewOrderProductUnsafe(
	orderID vObject.OrderID,
	productID vObject.ProductID,
	opts ...Option[*OrderProduct],
) OrderProduct {
	op := OrderProduct{
		OrderID:   orderID,
		ProductID: productID,
	}

	for _, opt := range opts {
		_ = opt(&op)
	}

	tn := op.Now()
	op.CreatedAt = tn
	op.UpdatedAt = tn

	return op
}

package addproducttoorder

import "github.com/google/uuid"

type Requestable interface {
	GetOrderID() uuid.UUID
	GetUserID() uuid.UUID
	GetProductID() uuid.UUID
	GetQuantity() uint64
}

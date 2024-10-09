package addproducttoorder

import "github.com/google/uuid"

type testRequest struct {
	orderUUID   uuid.UUID
	userUUID    uuid.UUID
	productUUID uuid.UUID
	quantity    uint64
}

var _ Requestable = (*testRequest)(nil)

func (t testRequest) GetOrderID() uuid.UUID {
	return t.orderUUID
}

func (t testRequest) GetUserID() uuid.UUID {
	return t.userUUID
}

func (t testRequest) GetProductID() uuid.UUID {
	return t.productUUID
}

func (t testRequest) GetQuantity() uint64 {
	return t.quantity
}

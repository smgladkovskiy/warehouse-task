package valueobjects

type OrderStatus string

const (
	OrderStatusCreated OrderStatus = "created"
	OrderStatusPaid    OrderStatus = "paid"
	OrderStatusOrdered OrderStatus = "ordered"
	OrderStatusShipped OrderStatus = "shipped"

	OrderStatusReceived OrderStatus = "received"
	OrderStatusReturned OrderStatus = "returned"
	OrderStatusCanceled OrderStatus = "canceled"
)

var orderFlow = map[OrderStatus][]OrderStatus{
	OrderStatusCreated:  {OrderStatusPaid, OrderStatusCanceled},
	OrderStatusPaid:     {OrderStatusOrdered, OrderStatusCanceled},
	OrderStatusOrdered:  {OrderStatusShipped, OrderStatusReceived, OrderStatusCanceled},
	OrderStatusShipped:  {OrderStatusReceived, OrderStatusReturned, OrderStatusCanceled},
	OrderStatusReceived: {OrderStatusReturned},
}

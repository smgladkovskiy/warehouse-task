package entities

import (
	"time"

	vObject "github.com/smgladkovskiy/warehouse-task/internal/service/entities/value_objects"
)

type ProductMovement struct {
	ID            vObject.ProductMovementID
	ProductID     vObject.ProductID
	WarehouseID   vObject.WarehouseID
	OperationType vObject.OperationType
	Quantity      vObject.Quantity
	Price         vObject.Price
	CreatedAt     time.Time
}

type ProductMovements []ProductMovement

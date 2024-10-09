package entities

import (
	"time"

	vObject "github.com/smgladkovskiy/warehouse-task/internal/service/entities/value_objects"
)

type Warehouse struct {
	ID   vObject.WarehouseID
	Name vObject.WarehouseName
	//Capacity vObject.WarehouseCapacity // считаем, что склад бесконечный
	CreatedAt time.Time
	UpdateAt  time.Time
	DeleteAt  *time.Time

	Stocks Stocks
}

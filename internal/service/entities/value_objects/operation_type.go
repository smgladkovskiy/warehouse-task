package valueobjects

type OperationType string

const (
	OperationTypeIncome   OperationType = "income"    // Поступление товаров на склад
	OperationTypeReserve  OperationType = "reserve"   // Резерв товаров для продажи
	OperationTypeSale     OperationType = "sale"      // Продажа товаров
	OperationTypeTransfer OperationType = "transfer"  // Перемещение товара между складами
	OperationTypeWriteOff OperationType = "write_off" // Списание товаров
)

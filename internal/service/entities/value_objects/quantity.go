package valueobjects

type Quantity uint64

const QuantityZero Quantity = 0

func NewQuantityUnsafe(quantity uint64) Quantity {
	return Quantity(quantity)
}

func (q Quantity) IsLessThan(quantity uint64) bool {
	return q.Uint64() < quantity
}

func (q Quantity) Uint64() uint64 {
	return uint64(q)
}

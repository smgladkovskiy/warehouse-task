package valueobjects

type Price int64

func (p Price) Multiply(quantity Quantity) Price {
	return p * Price(quantity)
}

func (p *Price) Subtract(price Price) {
	*p -= price
}

func (p *Price) Add(price Price) {
	*p += price
}

func NewPriceUnsafe(cents int) Price {
	return Price(cents)
}

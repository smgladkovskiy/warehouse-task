package valueobjects

type ProductTitle string

func NewProductTitleUnsafe(title string) ProductTitle {
	return ProductTitle(title)
}

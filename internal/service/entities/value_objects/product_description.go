package valueobjects

type ProductDescription string

func NewProductDescriptionUnsafe(desc string) ProductDescription {
	return ProductDescription(desc)
}

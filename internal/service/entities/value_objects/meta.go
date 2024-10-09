package valueobjects

type (
	MetaTotal    int64
	MetaPerPage  int64
	MetaPage     int64
	MetaLastPage int64
	Meta         struct {
		Total    MetaTotal
		PerPage  MetaPerPage
		Page     MetaPage
		LastPage MetaLastPage
	}
)

const DefaultPerPage = 10

func NewMetaTotal(total int) MetaTotal {
	return MetaTotal(total)
}

func NewMetaPerPage(perPage int) MetaPerPage {
	if perPage <= 0 {
		perPage = DefaultPerPage
	}
	return MetaPerPage(perPage)
}

func NewMetaPage(page int) MetaPage {
	if page <= 0 {
		page = 1
	}
	return MetaPage(page)
}

func NewMetaLastPage(total int, perPage int) MetaLastPage {
	return MetaLastPage((total + perPage - 1) / perPage)
}

func (t MetaTotal) Int64() int64 {
	return int64(t)
}

func (p MetaPerPage) Int64() int64 {
	return int64(p)
}

func (p MetaPage) Int64() int64 {
	return int64(p)
}

func (p MetaLastPage) Int64() int64 {
	return int64(p)
}

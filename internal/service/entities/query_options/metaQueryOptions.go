package queryoptions

import vObject "github.com/smgladkovskiy/warehouse-task/internal/service/entities/value_objects"

type MetaQueryOptionable interface {
	ForLimit() uint64
	ForOffset() uint64
	GetMeta() *vObject.Meta
	WithMetaTotal(total int)
	setPage(page int)
	setPerPage(perPage int)
}

type MetaQueryOptions struct {
	Page    int
	PerPage int
	Total   int
}

var _ MetaQueryOptionable = (*MetaQueryOptions)(nil)

func NewMetaQueryOptions(opts ...QueryOption[*MetaQueryOptions]) *MetaQueryOptions {
	mqos := MetaQueryOptions{}

	for _, opt := range opts {
		opt(&mqos)
	}

	if mqos.PerPage == 0 {
		mqos.PerPage = vObject.DefaultPerPage
	}

	if mqos.Page == 0 {
		mqos.Page = 1
	}

	return &mqos
}

func WithMetaPage[T MetaQueryOptionable](page int) QueryOption[T] {
	return func(options T) {
		options.setPage(page)
	}
}

func WithMetaPerPage[T MetaQueryOptionable](perPage int) QueryOption[T] {
	return func(options T) {
		options.setPerPage(perPage)
	}
}

func (m *MetaQueryOptions) ForLimit() uint64 {
	return uint64(m.PerPage)
}

func (m *MetaQueryOptions) ForOffset() uint64 {
	if m.Page == 0 {
		m.Page = 1
	}

	return (uint64(m.Page) - 1) * m.ForLimit()
}

func (m *MetaQueryOptions) WithMetaTotal(total int) {
	m.Total = total
}

func (m MetaQueryOptions) GetMeta() *vObject.Meta {
	return &vObject.Meta{
		Total:    vObject.NewMetaTotal(m.Total),
		PerPage:  vObject.NewMetaPerPage(int(m.ForLimit())),
		Page:     vObject.NewMetaPage(m.Page),
		LastPage: vObject.NewMetaLastPage(m.Total, int(m.ForLimit())),
	}
}

func (m *MetaQueryOptions) setPage(page int) {
	m.Page = page
}

func (m *MetaQueryOptions) setPerPage(perPage int) {
	m.PerPage = perPage
}

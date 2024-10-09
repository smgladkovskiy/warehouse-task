package queryoptions

type QueryOptionable interface {
	IsForUpdate() bool
	setForUpdate()
}

type BasicQueryOptions struct {
	forUpdate bool
	fromSync  bool
}

var _ QueryOptionable = (*BasicQueryOptions)(nil)

type QueryOption[T any] func(options T)

func NewBasicQueryOptions(opts ...QueryOption[*BasicQueryOptions]) *BasicQueryOptions {
	bqos := BasicQueryOptions{}

	for _, opt := range opts {
		opt(&bqos)
	}

	return &bqos
}

func WithForUpdate[T QueryOptionable]() QueryOption[T] {
	return func(options T) {
		options.setForUpdate()
	}
}

func (s *BasicQueryOptions) setForUpdate() {
	s.forUpdate = true
}

func (s BasicQueryOptions) IsForUpdate() bool {
	return s.forUpdate
}

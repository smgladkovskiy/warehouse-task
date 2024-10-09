package queryoptions_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	queryoptions "github.com/smgladkovskiy/warehouse-task/internal/service/entities/query_options"
	vObject "github.com/smgladkovskiy/warehouse-task/internal/service/entities/value_objects"
)

func TestMetaQueryOptions_ForLimit(t *testing.T) {
	t.Parallel()

	ass := assert.New(t)

	metaQos := queryoptions.NewMetaQueryOptions()

	ass.Equal(uint64(vObject.DefaultPerPage), metaQos.ForLimit())
	ass.Equal(uint64(0), metaQos.ForOffset())

	metaQos = queryoptions.NewMetaQueryOptions(
		queryoptions.WithMetaPage[*queryoptions.MetaQueryOptions](4),
		queryoptions.WithMetaPerPage[*queryoptions.MetaQueryOptions](15),
	)

	metaQos.WithMetaTotal(237)

	ass.Equal(uint64(15), metaQos.ForLimit())
	ass.Equal(uint64(45), metaQos.ForOffset())
	ass.Equal(vObject.Meta{
		Total:    237,
		PerPage:  15,
		Page:     4,
		LastPage: 16,
	}, *metaQos.GetMeta())
}

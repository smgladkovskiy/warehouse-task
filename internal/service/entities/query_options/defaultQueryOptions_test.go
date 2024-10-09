package queryoptions_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	queryoptions "github.com/smgladkovskiy/warehouse-task/internal/service/entities/query_options"
)

func TestBasicQueryOptions_IsForUpdate(t *testing.T) {
	t.Parallel()

	ass := assert.New(t)

	ass.True(queryoptions.NewBasicQueryOptions(queryoptions.WithForUpdate[*queryoptions.BasicQueryOptions]()).IsForUpdate())
	ass.False(queryoptions.NewBasicQueryOptions().IsForUpdate())
}

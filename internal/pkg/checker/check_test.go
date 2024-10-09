//go:build unit

package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/checker"
)

type handler struct{}

type testUseCase struct {
	checker.WithCheck

	requiredHandler *handler
	optionalHandler *handler `check:"optional"`
}

func TestWithCheck_Check(t *testing.T) {
	t.Parallel()

	uc := testUseCase{}

	require.ErrorIs(t, uc.Check(uc), checker.ErrInitError)

	uc.requiredHandler = &handler{}

	require.NoError(t, uc.Check(uc))
}

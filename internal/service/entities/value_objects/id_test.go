//go:build unit

package valueobjects_test

import (
	"database/sql/driver"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	vObject "github.com/smgladkovskiy/warehouse-task/internal/service/entities/value_objects"
)

//nolint:gochecknoglobals // сделано чисто для тестов
var ids = []vObject.IDInterface{}

func TestIDsName(t *testing.T) {
	t.Parallel()

	for _, id := range ids {
		tid := id.Clone()
		t.Run(tid.Name(), func(t *testing.T) {
			t.Parallel()

			expName := strings.ReplaceAll(reflect.TypeOf(tid).String(), "valueobject.", "")
			assert.Equal(t, expName, "*"+tid.Name())
		})
	}
}

func TestIDsSetEmpty(t *testing.T) {
	t.Parallel()

	for _, id := range ids {
		tid := id.Clone()

		t.Run(tid.Name(), func(t *testing.T) {
			t.Parallel()

			tid.SetFromUUID(uuid.New())

			tid.SetEmpty()

			assert.Equal(t, uuid.Nil, tid.UUID())
		})
	}
}

func TestIDsSetFromString(t *testing.T) {
	t.Parallel()

	for _, id := range ids {
		idSetFromString(t, id.Clone())
	}
}

func TestIDsString(t *testing.T) {
	t.Parallel()

	for _, id := range ids {
		idString(t, id.Clone())
	}
}

func TestIDsBytes(t *testing.T) {
	t.Parallel()

	for _, id := range ids {
		idBytes(t, id.Clone())
	}
}

func TestIDsUUID(t *testing.T) {
	t.Parallel()

	for _, id := range ids {
		idUUID(t, id.Clone())
	}
}

func TestIDsIsNull(t *testing.T) {
	t.Parallel()

	for _, id := range ids {
		idIsNull(t, id.Clone())
	}
}

func TestIDsMarshalJSON(t *testing.T) {
	t.Parallel()

	for _, id := range ids {
		idMarshalJSON(t, id.Clone())
	}
}

func TestIDsUnmarshalJSON(t *testing.T) {
	t.Parallel()

	for _, id := range ids {
		idUnmarshalJSON(t, id.Clone())
	}
}

func TestIDsScan(t *testing.T) {
	t.Parallel()

	for _, id := range ids {
		idScan(t, id.Clone())
	}
}

func TestIDsValue(t *testing.T) {
	t.Parallel()

	for _, id := range ids {
		idValue(t, id.Clone())
	}
}

func TestEqualsAnother(t *testing.T) {
	t.Parallel()

	for _, id := range ids {
		idEqualsAnother(t, id.Clone())
	}
}

func idSetFromString(t *testing.T, id vObject.IDInterface) {
	t.Helper()

	type tc struct {
		name   string
		in     string
		exp    uuid.UUID
		hasErr bool
	}

	tcs := []tc{
		{
			name: "правильный uuid",
			in:   "93786102-93CE-4F83-BD01-B93168695563",
			exp:  uuid.MustParse("93786102-93CE-4F83-BD01-B93168695563"),
		},
		{
			name: "пустой uuid",
			in:   "",
			exp:  uuid.Nil,
		},
		{
			name: "нулёный uuid",
			in:   vObject.NilUUIDStr,
			exp:  uuid.Nil,
		},
		{
			name:   "неправильный uuid",
			in:     "93786102-93CE-4F83-BD01-B9316869556",
			hasErr: true,
		},
	}

	for _, tc := range tcs {
		tc := tc
		tid := id.Clone()
		t.Run(tid.Name()+" "+tc.name, func(t *testing.T) {
			t.Parallel()

			err := tid.SetFromString(tc.in)
			if tc.hasErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.exp, tid.UUID())
		})
	}
}

func idString(t *testing.T, id vObject.IDInterface) {
	t.Helper()

	type tc struct {
		name string
		in   string
		exp  string
	}

	tcs := []tc{
		{
			name: "normal uuid",
			in:   "d241e000-3a6b-4afa-89ec-b17689577fdc",
			exp:  "d241e000-3a6b-4afa-89ec-b17689577fdc",
		},
		{
			name: "empty uuid",
			in:   "",
			exp:  uuid.Nil.String(),
		},
	}

	for _, tc := range tcs {
		tc := tc
		tid := id.Clone()
		t.Run(tid.Name()+" "+tc.name, func(t *testing.T) {
			t.Parallel()

			_ = tid.SetFromString(tc.in)

			out := tid.String()

			assert.Equal(t, tc.exp, out)
		})
	}
}

func idBytes(t *testing.T, id vObject.IDInterface) {
	t.Helper()

	type tc struct {
		name string
		in   string
		exp  []byte
	}

	tcs := []tc{
		{
			name: "normal uuid",
			in:   "d241e000-3a6b-4afa-89ec-b17689577fdc",
			exp:  []byte(`d241e000-3a6b-4afa-89ec-b17689577fdc`),
		},
		{
			name: "empty uuid",
			in:   vObject.NilUUIDStr,
			exp:  []byte(vObject.NilUUIDStr),
		},
		{
			name: "nil uuid",
			in:   "",
			exp:  []byte(vObject.NilUUIDStr),
		},
	}

	for _, tc := range tcs {
		tc := tc
		tid := id.Clone()
		t.Run(tid.Name()+" "+tc.name, func(t *testing.T) {
			t.Parallel()

			if err := tid.SetFromString(tc.in); !assert.NoError(t, err) {
				t.FailNow()
			}

			out := tid.Bytes()

			assert.Equal(t, tc.exp, out)
		})
	}
}

func idUUID(t *testing.T, id vObject.IDInterface) {
	t.Helper()

	type tc struct {
		name string
		in   string
		exp  uuid.UUID
	}

	tcs := []tc{
		{
			name: "normal uuid",
			in:   "d241e000-3a6b-4afa-89ec-b17689577fdc",
			exp:  uuid.MustParse("d241e000-3a6b-4afa-89ec-b17689577fdc"),
		},
		{
			name: "empty uuid",
			in:   uuid.Nil.String(),
			exp:  uuid.Nil,
		},
	}

	for _, tc := range tcs {
		tc := tc
		tid := id.Clone()
		t.Run(tid.Name()+" "+tc.name, func(t *testing.T) {
			t.Parallel()

			if err := tid.SetFromString(tc.in); !assert.NoError(t, err) {
				t.FailNow()
			}

			out := tid.UUID()

			assert.Equal(t, tc.exp, out)
		})
	}
}

func idIsNull(t *testing.T, id vObject.IDInterface) {
	t.Helper()

	type tc struct {
		name string
		in   string
		exp  bool
	}

	tcs := []tc{
		{
			name: "null value",
			in:   "",
			exp:  true,
		},
		{
			name: "zero value",
			in:   vObject.NilUUIDStr,
			exp:  true,
		},
		{
			name: "not null value",
			in:   uuid.New().String(),
			exp:  false,
		},
	}

	for _, tc := range tcs {
		tc := tc
		tid := id.Clone()

		t.Run(tid.Name()+" "+tc.name, func(t *testing.T) {
			t.Parallel()

			err := tid.SetFromString(tc.in)
			require.NoError(t, err)

			assert.Equal(t, tc.exp, tid.IsNil())
		})
	}
}

func idMarshalJSON(t *testing.T, id vObject.IDInterface) {
	t.Helper()

	type tc struct {
		name string
		in   string
		exp  []byte
		err  bool
	}

	nilIDBytes, _ := uuid.Nil.MarshalText()
	nilID := make([]byte, 0)
	nilID = append(nilID, `"`...)
	nilID = append(nilID, nilIDBytes...)
	nilID = append(nilID, `"`...)
	uuidID := uuid.New()
	uuidIDBytes, _ := uuidID.MarshalText()

	uuidIDBytes = append(uuidIDBytes, `"`...)
	uuidIDBytes = append([]byte(`"`), uuidIDBytes...)

	tcs := []tc{
		{
			name: "normal id",
			in:   uuidID.String(),
			exp:  uuidIDBytes,
		},
		{
			name: "zero id",
			in:   vObject.NilUUIDStr,
			exp:  nilID,
		},
		{
			name: "null id",
			in:   "",
			exp:  nilID,
		},
	}

	for _, tc := range tcs {
		tc := tc
		tid := id.Clone()
		t.Run(tid.Name()+" "+tc.name, func(t *testing.T) {
			t.Parallel()

			err := tid.SetFromString(tc.in)
			require.NoError(t, err)

			out, err := tid.MarshalJSON()
			if tc.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, string(tc.exp), string(out))
		})
	}
}

func idUnmarshalJSON(t *testing.T, id vObject.IDInterface) {
	t.Helper()

	type tc struct {
		name string
		in   []byte
		exp  uuid.UUID
	}

	uuidNew := uuid.New()
	idBB, _ := uuidNew.MarshalText()
	idBB = append(idBB, `"`...)
	idBB = append([]byte(`"`), idBB...)

	zeroIDBB, _ := uuid.Nil.MarshalText()

	zeroIDBB = append(zeroIDBB, `"`...)
	zeroIDBB = append([]byte(`"`), zeroIDBB...)

	tcs := []tc{
		{
			name: "корректный uuid",
			in:   idBB,
			exp:  uuidNew,
		},
		{
			name: "nil uuid",
			in:   nil,
			exp:  uuid.Nil,
		},
		{
			name: "zero uuid",
			in:   zeroIDBB,
			exp:  uuid.Nil,
		},
	}

	for _, tc := range tcs {
		tc := tc
		tid := id.Clone()

		t.Run(tid.Name()+" "+tc.name, func(t *testing.T) {
			t.Parallel()

			if err := tid.UnmarshalJSON(tc.in); !assert.NoError(t, err) {
				t.FailNow()
			}

			assert.Equal(t, tc.exp, tid.UUID())
		})
	}
}

func idScan(t *testing.T, id vObject.IDInterface) {
	t.Helper()

	type tc struct {
		name string
		in   interface{}
		exp  uuid.UUID
		err  error
	}

	idUUID := uuid.New()

	tcs := []tc{
		{
			name: "корректная строка UUID",
			in:   idUUID.String(),
			exp:  idUUID,
		},
		{
			name: "nil значение UUID",
			in:   nil,
			exp:  uuid.Nil,
		},
		{
			name: "zero значение UUID",
			in:   uuid.Nil.String(),
			exp:  uuid.Nil,
		},
		{
			name: "не строковое значение UUID",
			in:   struct{}{},
			exp:  uuid.Nil,
			err:  vObject.ErrScanError,
		},
	}

	for _, tc := range tcs {
		tc := tc
		tid := id.Clone()
		t.Run(tid.Name()+" "+tc.name, func(t *testing.T) {
			t.Parallel()

			err := tid.Scan(tc.in)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.exp, tid.UUID())
		})
	}
}

func idValue(t *testing.T, id vObject.IDInterface) {
	t.Helper()

	type tc struct {
		name string
		in   uuid.UUID
		exp  driver.Value
	}

	idUUID := uuid.New()

	tcs := []tc{
		{
			name: "корректный UUID",
			in:   idUUID,
			exp:  idUUID.String(),
		},
		{
			name: "nil значение UUID",
			in:   uuid.Nil,
			exp:  nil,
		},
	}

	for _, tc := range tcs {
		tc := tc
		tid := id.Clone()
		t.Run(tid.Name()+" "+tc.name, func(t *testing.T) {
			t.Parallel()

			tid.SetFromUUID(tc.in)

			out, err := tid.Value()
			if !assert.NoError(t, err) {
				t.FailNow()
			}

			assert.Equal(t, tc.exp, out)
		})
	}
}

func idEqualsAnother(t *testing.T, id vObject.IDInterface) {
	t.Helper()

	type tc struct {
		name string
		in   uuid.UUID
		exp  bool
	}

	idUUID := uuid.New()

	tcs := []tc{
		{
			name: "равные id",
			in:   idUUID,
			exp:  true,
		},
		{
			name: "отличающиеся id",
			in:   uuid.Nil,
			exp:  false,
		},
	}

	for _, tc := range tcs {
		tc := tc
		tid := id.Clone()
		anotherID := tid.Clone()
		t.Run(tid.Name()+" "+tc.name, func(t *testing.T) {
			t.Parallel()

			tid.SetFromUUID(tc.in)
			anotherID.SetFromUUID(idUUID)

			out := tid.EqualsAnother(anotherID)
			assert.Equal(t, tc.exp, out)
		})
	}
}

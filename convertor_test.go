package gmvc

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertor(t *testing.T) {

	cases := []struct {
		name   string
		origin string
		typ    any
		want   any
	}{

		// cases of non-pointer types
		{
			name:   "int",
			origin: "1",
			typ:    int(0),
			want:   int(1),
		},
		{
			name:   "int8",
			origin: "1",
			typ:    int8(0),
			want:   int8(1),
		},
		{
			name:   "int16",
			origin: "1",
			typ:    int16(0),
			want:   int16(1),
		},
		{
			name:   "int32",
			origin: "1",
			typ:    int32(0),
			want:   int32(1),
		},
		{
			name:   "int64",
			origin: "1",
			typ:    int64(0),
			want:   int64(1),
		},
		{
			name:   "uint",
			origin: "1",
			typ:    uint(0),
			want:   uint(1),
		},
		{
			name:   "uint8",
			origin: "1",
			typ:    uint8(0),
			want:   uint8(1),
		},
		{
			name:   "uint16",
			origin: "1",
			typ:    uint16(0),
			want:   uint16(1),
		},
		{
			name:   "uint32",
			origin: "1",
			typ:    uint32(0),
			want:   uint32(1),
		},
		{
			name:   "uint64",
			origin: "1",
			typ:    uint64(0),
			want:   uint64(1),
		},
		{
			name:   "float32",
			origin: "1.0",
			typ:    float32(0),
			want:   float32(1.0),
		},
		{
			name:   "float64",
			origin: "1.0",
			typ:    float64(0),
			want:   float64(1.0),
		},
		{
			name:   "bool",
			origin: "true",
			typ:    bool(false),
			want:   bool(true),
		},
		{
			name:   "string",
			origin: "1",
			typ:    "",
			want:   "1",
		},

		// cases of pointer types
		{
			name:   "int point",
			origin: "1",
			typ:    ptr(int(0)),
			want:   ptr(int(1)),
		},
		{
			name:   "int8 point",
			origin: "1",
			typ:    ptr(int8(0)),
			want:   ptr(int8(1)),
		},
		{
			name:   "int16 point",
			origin: "1",
			typ:    ptr(int16(0)),
			want:   ptr(int16(1)),
		},
		{
			name:   "int32 point",
			origin: "1",
			typ:    ptr(int32(0)),
			want:   ptr(int32(1)),
		},
		{
			name:   "int64 point",
			origin: "1",
			typ:    ptr(int64(0)),
			want:   ptr(int64(1)),
		},
		{
			name:   "uint point",
			origin: "1",
			typ:    ptr(uint(0)),
			want:   ptr(uint(1)),
		},
		{
			name:   "uint8 point",
			origin: "1",
			typ:    ptr(uint8(0)),
			want:   ptr(uint8(1)),
		},
		{
			name:   "uint16 point",
			origin: "1",
			typ:    ptr(uint16(0)),
			want:   ptr(uint16(1)),
		},
		{
			name:   "uint32 point",
			origin: "1",
			typ:    ptr(uint32(0)),
			want:   ptr(uint32(1)),
		},
		{
			name:   "uint64 point",
			origin: "1",
			typ:    ptr(uint64(0)),
			want:   ptr(uint64(1)),
		},
		{
			name:   "float32 point",
			origin: "1.0",
			typ:    ptr(float32(0)),
			want:   ptr(float32(1.0)),
		},
		{
			name:   "float64 point",
			origin: "1.0",
			typ:    ptr(float64(0)),
			want:   ptr(float64(1.0)),
		},
		{
			name:   "bool point",
			origin: "true",
			typ:    ptr(bool(false)),
			want:   ptr(bool(true)),
		},
		{
			name:   "string point",
			origin: "1",
			typ:    ptr(""),
			want:   ptr("1"),
		},

		// cases of slice types
		{
			name:   "int slice",
			origin: "1,2,3",
			typ:    []int{},
			want:   []int{1, 2, 3},
		},
		{
			name:   "int8 slice",
			origin: "1,2,3",
			typ:    []int8{},
			want:   []int8{1, 2, 3},
		},
		{
			name:   "int16 slice",
			origin: "1,2,3",
			typ:    []int16{},
			want:   []int16{1, 2, 3},
		},
		{
			name:   "int32 slice",
			origin: "1,2,3",
			typ:    []int32{},
			want:   []int32{1, 2, 3},
		},
		{
			name:   "int64 slice",
			origin: "1,2,3",
			typ:    []int64{},
			want:   []int64{1, 2, 3},
		},
		{
			name:   "uint slice",
			origin: "1,2,3",
			typ:    []uint{},
			want:   []uint{1, 2, 3},
		},
		{
			name:   "uint8 slice",
			origin: "1,2,3",
			typ:    []uint8{},
			want:   []uint8{1, 2, 3},
		},
		{
			name:   "uint16 slice",
			origin: "1,2,3",
			typ:    []uint16{},
			want:   []uint16{1, 2, 3},
		},
		{
			name:   "uint32 slice",
			origin: "1,2,3",
			typ:    []uint32{},
			want:   []uint32{1, 2, 3},
		},
		{
			name:   "uint64 slice",
			origin: "1,2,3",
			typ:    []uint64{},
			want:   []uint64{1, 2, 3},
		},
		{
			name:   "float32 slice",
			origin: "1.0,2.0,3.0",
			typ:    []float32{},
			want:   []float32{1.0, 2.0, 3.0},
		},
		{
			name:   "float64 slice",
			origin: "1.0,2.0,3.0",
			typ:    []float64{},
			want:   []float64{1.0, 2.0, 3.0},
		},
		{
			name:   "bool slice",
			origin: "true,false,true",
			typ:    []bool{},
			want:   []bool{true, false, true},
		},
		{
			name:   "string slice",
			origin: "1,2,3",
			typ:    []string{},
			want:   []string{"1", "2", "3"},
		},

		// cases of slice of pointer types
		{
			name:   "int slice of pointer",
			origin: "1,2,3",
			typ:    []*int{},
			want:   []*int{ptr(int(1)), ptr(int(2)), ptr(int(3))},
		},
		{
			name:   "int8 slice of pointer",
			origin: "1,2,3",
			typ:    []*int8{},
			want:   []*int8{ptr(int8(1)), ptr(int8(2)), ptr(int8(3))},
		},
		{
			name:   "int16 slice of pointer",
			origin: "1,2,3",
			typ:    []*int16{},
			want:   []*int16{ptr(int16(1)), ptr(int16(2)), ptr(int16(3))},
		},
		{
			name:   "int32 slice of pointer",
			origin: "1,2,3",
			typ:    []*int32{},
			want:   []*int32{ptr(int32(1)), ptr(int32(2)), ptr(int32(3))},
		},
		{
			name:   "int64 slice of pointer",
			origin: "1,2,3",
			typ:    []*int64{},
			want:   []*int64{ptr(int64(1)), ptr(int64(2)), ptr(int64(3))},
		},
		{
			name:   "uint slice of pointer",
			origin: "1,2,3",
			typ:    []*uint{},
			want:   []*uint{ptr(uint(1)), ptr(uint(2)), ptr(uint(3))},
		},
		{
			name:   "uint8 slice of pointer",
			origin: "1,2,3",
			typ:    []*uint8{},
			want:   []*uint8{ptr(uint8(1)), ptr(uint8(2)), ptr(uint8(3))},
		},
		{
			name:   "uint16 slice of pointer",
			origin: "1,2,3",
			typ:    []*uint16{},
			want:   []*uint16{ptr(uint16(1)), ptr(uint16(2)), ptr(uint16(3))},
		},
		{
			name:   "uint32 slice of pointer",
			origin: "1,2,3",
			typ:    []*uint32{},
			want:   []*uint32{ptr(uint32(1)), ptr(uint32(2)), ptr(uint32(3))},
		},
		{
			name:   "uint64 slice of pointer",
			origin: "1,2,3",
			typ:    []*uint64{},
			want:   []*uint64{ptr(uint64(1)), ptr(uint64(2)), ptr(uint64(3))},
		},
		{
			name:   "float32 slice of pointer",
			origin: "1.0,2.0,3.0",
			typ:    []*float32{},
			want:   []*float32{ptr(float32(1.0)), ptr(float32(2.0)), ptr(float32(3.0))},
		},
		{
			name:   "float64 slice of pointer",
			origin: "1.0,2.0,3.0",
			typ:    []*float64{},
			want:   []*float64{ptr(float64(1.0)), ptr(float64(2.0)), ptr(float64(3.0))},
		},
		{
			name:   "bool slice of pointer",
			origin: "true,false,true",
			typ:    []*bool{},
			want:   []*bool{ptr(bool(true)), ptr(bool(false)), ptr(bool(true))},
		},
		{
			name:   "string slice of pointer",
			origin: "1,2,3",
			typ:    []*string{},
			want:   []*string{ptr("1"), ptr("2"), ptr("3")},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ret, err := Convert(c.origin, reflect.TypeOf(c.typ))
			assert.Nil(t, err)
			assert.Equal(t, c.want, ret)
		})
	}

}

func ptr[T any](v T) *T {
	return &v
}

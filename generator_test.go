package goloop

import (
	"fmt"
	"reflect"
	"testing"

	"golang.org/x/exp/constraints"
)

type intGenTest[T constraints.Integer] struct {
	start, stop, step T
	want              []T
	wantPanic         bool
}

func testIntGenHelper[T constraints.Integer](t *testing.T, tests []intGenTest[T]) {
	t.Helper()

	var (
		foo  T
		kind = reflect.TypeOf(foo).Kind()
	)
	for i, test := range tests {
		t.Run(fmt.Sprintf("%s %d", kind, i), func(t *testing.T) {
			t.Helper()

			defer func() {
				e := recover()
				if test.wantPanic && e == nil {
					t.Errorf("test case %v, want panic, didn't panic", test)
				}
				if !test.wantPanic && e != nil {
					t.Errorf("test case %v, unexpected panic", test)
				}
			}()

			var (
				got []T
				gen = newIntGen(test.start, test.stop, test.step)
			)
			for gen.next() {
				got = append(got, gen.gen())
			}
			if !reflect.DeepEqual(test.want, got) {
				t.Errorf("(*intGen).gen() generates %v, want %v", got, test.want)
			}
		})
	}
}

func TestIntGen(t *testing.T) {
	testIntGenHelper(t, []intGenTest[int8]{
		{0, 0, 1, nil, true},
		{0, 1, 0, nil, true},
		{0, 1, -1, nil, true},
		{-3, 3, 1, []int8{-3, -2, -1, 0, 1, 2, 3}, false},
		{-3, 3, 2, []int8{-3, -1, 1, 3}, false},
		{-3, 3, 3, []int8{-3, 0, 3}, false},
		{4, -2, 1, []int8{4, 3, 2, 1, 0, -1, -2}, false},
		{4, -3, 2, []int8{4, 2, 0, -2}, false},
		{4, -1, 3, []int8{4, 1}, false},
		{120, 127, 1, []int8{120, 121, 122, 123, 124, 125, 126, 127}, false},
		{120, 127, 3, []int8{120, 123, 126}, false},
		{-120, -128, 1, []int8{-120, -121, -122, -123, -124, -125, -126, -127, -128}, false},
		{-120, -128, 3, []int8{-120, -123, -126}, false},
	})
	testIntGenHelper(t, []intGenTest[uint8]{
		{0, 0, 1, nil, true},
		{4, 4, 1, nil, true},
		{1, 4, 0, nil, true},
		{0, 5, 1, []uint8{0, 1, 2, 3, 4, 5}, false},
		{0, 5, 2, []uint8{0, 2, 4}, false},
		{1, 5, 1, []uint8{1, 2, 3, 4, 5}, false},
		{1, 5, 3, []uint8{1, 4}, false},
		{250, 255, 1, []uint8{250, 251, 252, 253, 254, 255}, false},
		{250, 255, 2, []uint8{250, 252, 254}, false},
		{250, 255, 3, []uint8{250, 253}, false},
	})
	testIntGenHelper(t, []intGenTest[int]{
		{-4, -4, 1, nil, true},
		{1, 4, 0, nil, true},
		{1, 4, -1, nil, true},
		{-1, 3, 1, []int{-1, 0, 1, 2, 3}, false},
		{2, -3, 1, []int{2, 1, 0, -1, -2, -3}, false},
	})
	testIntGenHelper(t, []intGenTest[uint]{
		{0, 0, 1, nil, true},
		{1, 4, 0, nil, true},
		{10, 23, 4, []uint{10, 14, 18, 22}, false},
		{10, 3, 4, []uint{10, 6}, false},
	})
}

func testIntGenOneHelper[T constraints.Integer](t *testing.T, value T) {
	t.Helper()

	var (
		got []T
		gen = newIntGenOne(value)
	)
	for gen.next() {
		got = append(got, gen.gen())
	}
	if !reflect.DeepEqual(got, []T{value}) {
		t.Errorf("(*intGenOne).gen() generates %v, want %v", got, []T{value})
	}
}

func TestIntGenOne(t *testing.T) {
	testIntGenOneHelper(t, int8(127))
	testIntGenOneHelper(t, int8(-128))
	testIntGenOneHelper(t, uint8(0))
	testIntGenOneHelper(t, uint8(255))
	testIntGenOneHelper(t, int(399))
	testIntGenOneHelper(t, int(-399))
	testIntGenOneHelper(t, int32(67))
	testIntGenOneHelper(t, uint64(67))
}

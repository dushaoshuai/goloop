package goloop_test

import (
	"fmt"
	"reflect"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
	"go.uber.org/goleak"
	"golang.org/x/exp/constraints"

	"github.com/dushaoshuai/goloop"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestRepeat(t *testing.T) {
	for _, test := range []struct {
		times int
		want  []int
	}{
		{times: -1, want: nil},
		{times: 0, want: nil},
		{times: 1, want: []int{0}},
		{times: 2, want: []int{0, 1}},
	} {
		t.Run(fmt.Sprintf("%d times", test.times), func(t *testing.T) {
			var got []int
			for i := range goloop.Repeat(test.times) {
				got = append(got, i)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Wrong iteration values produced by goloop.Repeat, want %v, got %v", test.want, got)
			}
		})
	}
}

func TestRepeatWithBreak(t *testing.T) {
	for _, test := range []struct {
		times      int
		breakpoint int
		want       []int
	}{
		{times: -1, breakpoint: 0, want: nil},
		{times: 0, breakpoint: 0, want: nil},
		{times: 1, breakpoint: 0, want: []int{0}},
		{times: 2, breakpoint: 1, want: []int{0, 1}},
		{times: 2, breakpoint: 0, want: []int{0}},
		{times: 3, breakpoint: 2, want: []int{0, 1, 2}},
		{times: 3, breakpoint: 1, want: []int{0, 1}},
		{times: 3, breakpoint: 3, want: []int{0, 1, 2}},
		{times: 3, breakpoint: 4, want: []int{0, 1, 2}},
	} {
		t.Run(fmt.Sprintf("%d times", test.times), func(t *testing.T) {
			var got []int
			for i := range goloop.RepeatWithBreak(test.times) {
				got = append(got, i.I)
				if i.I == test.breakpoint {
					i.Break()
				}
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("RepeatWithBreak generates %v, want %v", got, test.want)
			}
		})
	}

}

type testRange[T constraints.Integer] struct {
	wantPanic   bool
	start, stop T
	step        []T
	want        []T
	breakpoints mapset.Set[T]
}

func testRangeHelper[T constraints.Integer](t *testing.T, tests []testRange[T]) {
	t.Helper()

	var (
		foo  T
		kind = reflect.TypeOf(foo).Kind()
	)

	for i, test := range tests {
		t.Run(fmt.Sprintf("%s-%d", kind, i), func(t *testing.T) {
			t.Helper()
			defer func() {
				e := recover()
				if test.wantPanic && e == nil {
					t.Errorf("test Range, want panic, but didn't happen")
				}
				if !test.wantPanic && e != nil {
					t.Errorf("test range, unexpected panic: %v", e)
				}
			}()

			var (
				got []T
			)
			for j := range goloop.Range(test.start, test.stop, test.step...) {
				got = append(got, j.I)
				if bps := test.breakpoints; bps != nil && bps.Contains(j.I) {
					j.Break()
				}
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Range() generates %v, want %v", got, test.want)
			}
		})
	}
}

func TestRange(t *testing.T) {
	testRangeHelper[int8](t, []testRange[int8]{
		{true, 0, 1, []int8{0}, nil, nil},
		{true, 0, 1, []int8{-1}, nil, nil},
		{false, 0, 0, nil, []int8{0}, nil},
		{false, 0, 0, []int8{1}, []int8{0}, nil},
		{false, 0, 1, nil, []int8{0, 1}, nil},
		{false, 0, 1, nil, []int8{0}, mapset.NewThreadUnsafeSet[int8](0)},
		{false, -10, 10, []int8{4}, []int8{-10, -6, -2, 2, 6, 10}, nil},
		{false, -10, 10, []int8{4}, []int8{-10, -6, -2}, mapset.NewThreadUnsafeSet[int8](-2, 2)},
		{false, 10, -10, []int8{6}, []int8{10, 4, -2, -8}, nil},
		{false, 10, -10, []int8{6}, []int8{10, 4, -2}, mapset.NewThreadUnsafeSet[int8](-2, -3)},
	})
	testRangeHelper[uint8](t, []testRange[uint8]{
		{true, 0, 1, []uint8{0}, nil, nil},
		{false, 0, 0, []uint8{1}, []uint8{0}, nil},
		{false, 0, 5, nil, []uint8{0, 1, 2, 3, 4, 5}, nil},
		{false, 0, 5, []uint8{2}, []uint8{0, 2, 4}, nil},
		{false, 0, 5, nil, []uint8{0, 1, 2, 3}, mapset.NewThreadUnsafeSet[uint8](3, 5)},
		{false, 0, 5, []uint8{2}, []uint8{0}, mapset.NewThreadUnsafeSet[uint8](0, 1)},
	})
}

type testRangeSlice[T constraints.Integer] struct {
	start, stop T
	step        []T
	want        []T
	wantPanic   bool
}

func testRangeSliceHelper[T constraints.Integer](t *testing.T, tests []testRangeSlice[T]) {
	t.Helper()

	var (
		foo  T
		kind = reflect.TypeOf(foo).Kind()
	)
	for i, test := range tests {
		t.Run(fmt.Sprintf("%s-%d", kind, i), func(t *testing.T) {
			defer func() {
				e := recover()
				if test.wantPanic && e == nil {
					t.Errorf("test RangeSlice, want panic, didn't happen")
				}
				if !test.wantPanic && e != nil {
					t.Errorf("test RangeSlice, unexpected panic")
				}
			}()

			got := goloop.RangeSlice(test.start, test.stop, test.step...)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("RangeSlice = %v, want %v", got, test.want)
			}
		})
	}
}

func TestRangeSlice(t *testing.T) {
	testRangeSliceHelper(t, []testRangeSlice[int8]{
		{1, 2, []int8{0}, nil, true},
		{1, 2, []int8{-1}, nil, true},
		{0, 0, nil, []int8{0}, false},
		{0, 0, []int8{1}, []int8{0}, false},
		{-2, 3, nil, []int8{-2, -1, 0, 1, 2, 3}, false},
		{5, -9, []int8{3}, []int8{5, 2, -1, -4, -7}, false},
	})
	testRangeSliceHelper(t, []testRangeSlice[uint8]{
		{1, 2, []uint8{0}, nil, true},
		{0, 0, nil, []uint8{0}, false},
		{0, 0, []uint8{1}, []uint8{0}, false},
		{0, 5, nil, []uint8{0, 1, 2, 3, 4, 5}, false},
		{0, 5, []uint8{1}, []uint8{0, 1, 2, 3, 4, 5}, false},
		{3, 10, []uint8{4}, []uint8{3, 7}, false},
	})
	testRangeSliceHelper(t, []testRangeSlice[int]{
		{1, 2, []int{0}, nil, true},
		{1, 2, []int{-1}, nil, true},
		{0, 0, nil, []int{0}, false},
		{0, 0, []int{1}, []int{0}, false},
		{-2, 3, nil, []int{-2, -1, 0, 1, 2, 3}, false},
		{5, -9, []int{3}, []int{5, 2, -1, -4, -7}, false},
	})
	testRangeSliceHelper(t, []testRangeSlice[uint]{
		{1, 2, []uint{0}, nil, true},
		{0, 0, nil, []uint{0}, false},
		{0, 0, []uint{1}, []uint{0}, false},
		{0, 5, []uint{1}, []uint{0, 1, 2, 3, 4, 5}, false},
		{3, 10, []uint{4}, []uint{3, 7}, false},
	})
}

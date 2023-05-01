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
				t.Errorf("Wrong iteration values produced by goloop.RepeatWithBreak, want %v, got %v", test.want, got)
			}
		})
	}

}

func TestRange(t *testing.T) {
	for _, test := range []struct {
		start, end int
		want       []int
	}{
		{-1, 1, []int{-1, 0, 1}},
		{1, -1, []int{1, 0, -1}},
		{1, 1, []int{1}},
		{-1, -1, []int{-1}},
		{-10, -7, []int{-10, -9, -8, -7}},
		{-7, -10, []int{-7, -8, -9, -10}},
		{3, 6, []int{3, 4, 5, 6}},
		{6, 3, []int{6, 5, 4, 3}},
	} {
		t.Run("", func(t *testing.T) {
			var got []int
			for i := range goloop.Range(test.start, test.end) {
				got = append(got, i.I)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Wrong iteration values produced by goloop.Range, want %v, got %v", test.want, got)
			}
		})
	}

}

func TestRangeBreak(t *testing.T) {
	for _, test := range []struct {
		start, end int
		breakpoint int
		want       []int
	}{
		{10, 10, 10, []int{10}},
		{-10, -10, -10, []int{-10}},
		{0, 0, 0, []int{0}},
		{-10, 10, -7, []int{-10, -9, -8, -7}},
		{10, -10, 6, []int{10, 9, 8, 7, 6}},
	} {
		t.Run("", func(t *testing.T) {
			var got []int
			for i := range goloop.Range(test.start, test.end) {
				got = append(got, i.I)
				if i.I == test.breakpoint {
					i.Break()
				}
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Wrong iteration values produced by goloop.Range, want %v, got %v", test.want, got)
			}
		})
	}
}

func TestRangeNoLeakWithMultiBreak(t *testing.T) {
	for _, test := range []struct {
		start, end  int
		breakpoints mapset.Set[int]
	}{
		{0, 0, mapset.NewThreadUnsafeSet(0)},
		{-10, 10, mapset.NewThreadUnsafeSet(7, 9, 0, -5, -4)},
		{10, -10, mapset.NewThreadUnsafeSet(7, 9, 0, -5, -4)},
	} {
		t.Run("", func(t *testing.T) {
			for i := range goloop.Range(test.start, test.end) {
				i := i
				go func() {
					if test.breakpoints.Contains(i.I) {
						i.Break()
					}
				}()
			}
		})
	}
}

func TestRangeWithStep(t *testing.T) {
	for _, test := range []struct {
		start, end, step int
		want             []int
		wantPanic        bool
	}{
		{0, 0, 1, []int{0}, false},
		{-1, 1, 1, []int{-1, 0, 1}, false},
		{-1, 1, -1, nil, true},
		{1, -1, 1, nil, true},
		{-4, 5, 2, []int{-4, -2, 0, 2, 4}, false},
		{-4, 5, -2, nil, true},
		{4, -3, -2, []int{4, 2, 0, -2}, false},
		{4, -3, 2, nil, true},
	} {
		t.Run("", func(t *testing.T) {
			if test.wantPanic {
				defer func() {
					err := recover()
					t.Log(err, test)
					if err == nil {
						t.Errorf("an inappropriate step was expected to cause panic, but it does not: %v", test)
					}
				}()
			}

			var got []int
			for i := range goloop.Range(test.start, test.end, uint64(test.step)) {
				got = append(got, i.I)
			}
			if !reflect.DeepEqual(test.want, got) {
				t.Errorf("Wrong iteration values produced by goloop.Range, want %v, got %v", test.want, got)
			}
		})
	}
}

type testRange[T constraints.Integer] struct {
	start, stop T
	step        []T
	want        []T
	wantPanic   bool
}

func testRangeSliceHelper[T constraints.Integer](t *testing.T, tests []testRange[T]) {
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
	testRangeSliceHelper(t, []testRange[int8]{
		{1, 2, []int8{0}, nil, true},
		{1, 2, []int8{-1}, nil, true},
		{0, 0, nil, []int8{0}, false},
		{0, 0, []int8{1}, []int8{0}, false},
		{-2, 3, nil, []int8{-2, -1, 0, 1, 2, 3}, false},
		{5, -9, []int8{3}, []int8{5, 2, -1, -4, -7}, false},
	})
	testRangeSliceHelper(t, []testRange[uint8]{
		{1, 2, []uint8{0}, nil, true},
		{0, 0, nil, []uint8{0}, false},
		{0, 0, []uint8{1}, []uint8{0}, false},
		{0, 5, nil, []uint8{0, 1, 2, 3, 4, 5}, false},
		{0, 5, []uint8{1}, []uint8{0, 1, 2, 3, 4, 5}, false},
		{3, 10, []uint8{4}, []uint8{3, 7}, false},
	})
	testRangeSliceHelper(t, []testRange[int]{
		{1, 2, []int{0}, nil, true},
		{1, 2, []int{-1}, nil, true},
		{0, 0, nil, []int{0}, false},
		{0, 0, []int{1}, []int{0}, false},
		{-2, 3, nil, []int{-2, -1, 0, 1, 2, 3}, false},
		{5, -9, []int{3}, []int{5, 2, -1, -4, -7}, false},
	})
	testRangeSliceHelper(t, []testRange[uint]{
		{1, 2, []uint{0}, nil, true},
		{0, 0, nil, []uint{0}, false},
		{0, 0, []uint{1}, []uint{0}, false},
		{0, 5, []uint{1}, []uint{0, 1, 2, 3, 4, 5}, false},
		{3, 10, []uint{4}, []uint{3, 7}, false},
	})
}

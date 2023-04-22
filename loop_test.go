package goloop_test

import (
	"fmt"
	"reflect"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
	"go.uber.org/goleak"

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

func ExampleRepeat() {
	for i := range goloop.Repeat(3) {
		fmt.Println("Repeat", i)
	}
	// Output:
	// Repeat 0
	// Repeat 1
	// Repeat 2
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

func ExampleRepeatWithBreak() {
	for i := range goloop.RepeatWithBreak(3) {
		fmt.Println("Repeat", i.I)
		if i.I == 1 {
			i.Break()
		}
	}
	// Output:
	// Repeat 0
	// Repeat 1
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
			for i := range goloop.Range(test.start, test.end, test.step) {
				got = append(got, i.I)
			}
			if !reflect.DeepEqual(test.want, got) {
				t.Errorf("Wrong iteration values produced by goloop.Range, want %v, got %v", test.want, got)
			}
		})
	}
}

func ExampleRange() {
	for i := range goloop.Range(3, 26, 5) {
		fmt.Println(i.I)
		if i.I == 18 {
			i.Break()
		}
	}

	// Output:
	// 3
	// 8
	// 13
	// 18
}

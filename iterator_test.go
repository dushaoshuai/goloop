package goloop

import (
	"reflect"
	"sync"
	"testing"

	"golang.org/x/exp/rand"
)

func TestChanIter(t *testing.T) {
	var ( // breakpoints
		bpn2 int8 = -2
		bpn1 int8 = -1
		bp0  int8 = 0
		bp1  int8 = 1
		bp2  int8 = 2
		bp3  int8 = 3
		bp9  int8 = 9
	)

	tests := []struct {
		values     []int8
		breakpoint *int8
		want       []int8
	}{
		{nil, nil, nil},
		{nil, &bp3, nil},
		{[]int8{-1, 0, 1, 2}, nil, []int8{-1, 0, 1, 2}},
		{[]int8{-1, 0, 1, 2}, &bpn2, []int8{-1, 0, 1, 2}},
		{[]int8{-1, 0, 1, 2}, &bpn1, []int8{-1}},
		{[]int8{-1, 0, 1, 2}, &bp0, []int8{-1, 0}},
		{[]int8{-1, 0, 1, 2}, &bp1, []int8{-1, 0, 1}},
		{[]int8{-1, 0, 1, 2}, &bp2, []int8{-1, 0, 1, 2}},
		{[]int8{-1, 0, 1, 2}, &bp3, []int8{-1, 0, 1, 2}},
		{[]int8{-1, 0, 1, 2}, &bp9, []int8{-1, 0, 1, 2}},
		{[]int8{1, 0, -1}, nil, []int8{1, 0, -1}},
		{[]int8{1, 0, -1}, &bp0, []int8{1, 0}},
		{[]int8{1, 0, -1}, &bpn1, []int8{1, 0, -1}},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			var (
				iter = newChanIter[int8]()
				wg   sync.WaitGroup
				got  []int8
			)

			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := range iter.c {
					got = append(got, i.I)
					if test.breakpoint != nil && *test.breakpoint == i.I {
						i.Break()
					}
				}
			}()

			for _, v := range test.values {
				if breaked := iter.iter(v); breaked {
					break
				}
			}
			iter.finish()

			wg.Wait()
			if !reflect.DeepEqual(test.want, got) {
				t.Errorf("chanIter generates %v, want %v", got, test.want)
			}
		})
	}
}

func TestChanIterNoLeakWithMultiBreak(t *testing.T) {
	for i := 0; i < 10000; i++ {
		t.Run("", func(t *testing.T) {
			var (
				iter = newChanIter[uint8]()
				wg   sync.WaitGroup
			)

			wg.Add(1)
			go func() {
				defer wg.Done()
				for _, j := range rand.Perm(256) {
					if breaked := iter.iter(uint8(j)); breaked {
						break
					}
				}
				iter.finish()
			}()

			for j := range iter.c {
				go func(i I[uint8]) {
					if i.I%2 == 0 {
						i.Break()
					}
				}(j)
			}
			wg.Wait()
		})
	}
}

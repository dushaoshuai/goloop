package goloop

import (
	"fmt"
)

// <= 0
// > 0
func Repeat(times int) <-chan int {
	c := make(chan int)
	go func() {
		for i := 0; i < times; i++ {
			c <- i
		}
		close(c)
	}()
	return c
}

func RepeatWithBreak(times int) <-chan I {
	rChan := Range(0, times-1)
	if times <= 0 {
		i := <-rChan
		i.Break()
	}
	return rChan
}

// As a special case, if start equals end, the iteration value produced is only start, no matter what the specified step is.
func Range(start, end int, step ...int) <-chan I {
	if start == end {
		return caseStartEqualsEnd(start)
	}

	var incr int
	if len(step) != 0 {
		incr = step[0]
		if (start < start+incr) != (start < end) {
			panic(fmt.Sprintf("goloop: infinite loop with start(%d)..end(%d)..step(%d)", start, end, incr))
		}
	} else {
		if start < end {
			incr = 1
		} else {
			incr = -1
		}
	}

	iter := newIterator()
	go func() {
		if start < end {
		L1:
			for i := start; i <= end; i += incr {
				if breaked := iter.iter(i); breaked {
					break L1
				}
			}
		} else {
		L2:
			for i := start; i >= end; i += incr {
				if breaked := iter.iter(i); breaked {
					break L2
				}
			}
		}
		iter.finish()
	}()
	return iter.c
}

func caseStartEqualsEnd(start int) <-chan I {
	iter := newIterator()
	go func() {
		iter.iter(start)
		iter.finish()
	}()
	return iter.c
}

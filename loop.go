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

type I struct {
	I     int
	Break func()
}

type iterator struct {
	c         chan I
	breakChan chan struct{}
}

func newIterator() *iterator {
	return &iterator{
		c:         make(chan I),
		breakChan: make(chan struct{}),
	}
}

func (i *iterator) breakFunc() {
	// Allow breakFunc to be called multiple times.
	defer func() {
		recover()
	}()

	close(i.breakChan)

	// Empty the channel to not immediately return to the caller of breakFunc
	// until the channel is closed. Prevent the caller receiving more iteration
	// values after called breakFunc.
	for range i.c {
	}
}

func (i *iterator) finish() {
	close(i.c)
}

func (i *iterator) iter(iterValue int) (breaked bool) {
	select {
	case <-i.breakChan:
		return true
	case i.c <- I{I: iterValue, Break: i.breakFunc}:
		return false
	}
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

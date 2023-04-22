// Package goloop tries to facilitate looping in Go.
// It imitates Go's "for ... range ... {}" looping style.
package goloop

import (
	"fmt"
)

// Repeat returns a read-only channel. Clients can iterate through values received
// on the returned channel to repeatedly doing something times times. No values
// will be sent on the channel if times is not greater than 0. Values will be
// sent in order and is in the range [0, times).
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

// RepeatWithBreak is almost the same as Repeat, except that the returned channel's
// element is I, whose Break field can be called to break the loop:
//
//	for i := range RepeatWithBreak(50) {
//		// Do something with i.I.
//		// Break the for loop if certain conditions are met.
//		if i == 30 {
//			i.Break()
//		}
//	}
func RepeatWithBreak(times int) <-chan I {
	rChan := Range(0, times-1)
	if times <= 0 {
		i := <-rChan
		i.Break()
	}
	return rChan
}

// Range returns a channel for the client to iterate. Values sent on the channel
// are start, start+step, start+2*step, ..., stop. If step is not specified, it
// defaults to 1 or -1 as appropriate. If the specified step causes an infinite
// loop, Range panics. As a special case, if start equals stop, the iteration
// value produced is only start, no matter what the specified step is.
func Range(start, stop int, step ...int) <-chan I {
	if start == stop {
		return caseStartEqualsEnd(start)
	}

	var incr int
	if len(step) != 0 {
		incr = step[0]
		if (start < start+incr) != (start < stop) {
			panic(fmt.Sprintf("goloop: infinite loop with start(%d)..stop(%d)..step(%d)", start, stop, incr))
		}
	} else {
		if start < stop {
			incr = 1
		} else {
			incr = -1
		}
	}

	iter := newIterator()
	go func() {
		if start < stop {
		L1:
			for i := start; i <= stop; i += incr {
				if breaked := iter.iter(i); breaked {
					break L1
				}
			}
		} else {
		L2:
			for i := start; i >= stop; i += incr {
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

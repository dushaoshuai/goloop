// Package goloop tries to facilitate looping in Go.
// It imitates Go's "for ... range ... {}" looping style.
package goloop

import (
	"golang.org/x/exp/constraints"
)

// Repeat returns a read-only channel. Clients can iterate through values received
// on the returned channel to repeatedly doing something times times.
// Values will be sent in order and is in the half-open interval [0,times).
// No values will be sent on the channel if times is less than or equal to 0.
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
//		if i.I == 30 {
//			i.Break()
//		}
//	}
func RepeatWithBreak(times int) <-chan I[int] {
	rChan := Range(0, times)
	if times <= 0 {
		i := <-rChan
		i.Break() // todo check
	}
	return rChan
}

// Range returns a read-only channel for the client to iterate. Values sent on
// the channel are start, start+step, start+2*step, ... with stop excluded.
// As a special case, if start equals stop, the iteration value produced is only
// start, no matter what the specified step is.
//
// todo step uint64
// If step is not specified, it defaults to 1.
// todo overflow
//
// The returned channel's element is I, whose Break field can be called to break the loop.
func Range[T constraints.Integer](start, stop T, step ...uint64) <-chan I[T] {
	if start == stop {
		return iterOnce(start)
	}

	var incr T
	if len(step) != 0 {
		incr = T(step[0])
	} else {
		incr = 1
	}

	iter := newIterator[T]()
	go func() {
		if start < stop {
			for i := start; i < stop; i += incr {
				if breaked := iter.iter(i); breaked {
					break
				}
				if i+incr < i { // overflow
					break
				}
			}
		} else {
			for i := start; i > stop; i -= incr {
				if breaked := iter.iter(i); breaked {
					break
				}
				if i-incr > i { // overflow
					break
				}
			}
		}
		iter.finish()
	}()
	return iter.c
}

func iterOnce[T constraint](value T) <-chan I[T] {
	iter := newIterator[T]()
	go func() {
		iter.iter(value)
		iter.finish()
	}()
	return iter.c
}

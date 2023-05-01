// Package goloop tries to facilitate looping in Go.
package goloop

import (
	"golang.org/x/exp/constraints"
)

// Repeat is intended to facilitate repeatedly doing something times times.
// Repeat generates a sequence of ints and send them on the returned channel.
// Values will be sent in order and are in the half-open interval [0,times).
// No values will be sent if times is less than or equal to 0.
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

// Range generates a sequence of integers and send them on the returned channel.
//
// If start is less than stop, the generated values are determined by the formula
// s[i] = start + step*i where s[i] is less than or equal to stop.
// If start is greater than stop, the generated values are determined by the formula
// s[i] = start - step*i where s[i] is greater than or equal to stop.
// If start is equal to stop, the only generated value is start(stop),
// no matter what the specified step is.
//
// If not specified, step is 1 by default. If specified, step must be greater than 0,
// otherwise Range will panic. There is one exception: if start equals stop,
// Range does not panic and generates one value: start(stop).
//
// The returned channel's element is I, whose Break field can be called to terminate communication.
func Range[T constraints.Integer](start, stop T, step ...T) <-chan I[T] {
	var gen generator[T]
	if start == stop {
		gen = newIntGenOne(start)
	} else {
		var incr T
		if len(step) != 0 {
			incr = step[0]
		} else {
			incr = 1
		}
		gen = newIntGen(start, stop, incr)
	}

	iter := newChanIter[T]()
	go func() {
		for gen.next() {
			if breaked := iter.iter(gen.gen()); breaked {
				break
			}
		}
		iter.finish()
	}()
	return iter.c
}

// RangeSlice generates a sequence of integers and put them in the returned slice.
//
// If start is less than stop, the generated values are determined by the formula
// s[i] = start + step*i where s[i] is less than or equal to stop.
// If start is greater than stop, the generated values are determined by the formula
// s[i] = start - step*i where s[i] is greater than or equal to stop.
// If start is equal to stop, the only generated value is start(stop),
// no matter what the specified step is.
//
// If not specified, step is 1 by default. If specified, step must be greater than 0,
// otherwise RangeSlice will panic. There is one exception: if start equals stop,
// RangeSlice does not panic and generates one value: start(stop).
func RangeSlice[T constraints.Integer](start, stop T, step ...T) (s []T) {
	var gen generator[T]
	if start == stop {
		gen = newIntGenOne(start)
	} else {
		var incr T
		if len(step) != 0 {
			incr = step[0]
		} else {
			incr = 1
		}
		gen = newIntGen(start, stop, incr)
	}

	var iter sliceIter[T]
	for gen.next() {
		iter = append(iter, gen.gen())
	}
	return iter
}

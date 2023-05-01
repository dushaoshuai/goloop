package goloop

import (
	"golang.org/x/exp/constraints"
)

// all types this package supported
type constraint interface {
	constraints.Integer
}

type generator[T constraint] interface {
	// next reports whether there is next T.
	// It returns false when the generator is exhausted.
	next() bool
	// gen generates the next T.
	gen() T
}

type intGen[T constraints.Integer] struct {
	started bool
	start   T // start must not be equal to stop
	stop    T
	step    T // step must be greater than 0
	curr    T
}

func newIntGen[T constraints.Integer](start, stop, step T) *intGen[T] {
	if start == stop {
		panic("goloop: intGen: use intGenOne instead if start is equal to stop")
	}
	if step <= 0 {
		panic("goloop: intGen: step must be greater than 0")
	}
	return &intGen[T]{
		start: start,
		stop:  stop,
		step:  step,
	}
}

func (g *intGen[T]) next() bool {
	if !g.started {
		g.started = true
		g.curr = g.start
		return true
	}
	if g.start < g.stop {
		if g.curr+g.step < g.curr { // check overflow
			return false
		}
		if g.curr+g.step <= g.stop {
			g.curr += g.step
			return true
		}
		return false
	} else {
		if g.curr-g.step > g.curr { // check overflow
			return false
		}
		if g.curr-g.step >= g.stop {
			g.curr -= g.step
			return true
		}
		return false
	}
}

func (g *intGen[T]) gen() T {
	return g.curr
}

// intGenOne generates only one value and is used when start equals stop.
type intGenOne[T constraints.Integer] struct {
	done  bool
	value T
}

func newIntGenOne[T constraints.Integer](value T) *intGenOne[T] {
	return &intGenOne[T]{
		done:  false,
		value: value,
	}
}

func (g *intGenOne[T]) next() bool {
	if !g.done {
		g.done = true
		return true
	}
	return false
}

func (g *intGenOne[T]) gen() T {
	return g.value
}

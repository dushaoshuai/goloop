package goloop

import (
	"fmt"
	"math"
	"reflect"

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
	stop    T // stop is excluded from the generated Ts
	step    T // step must be greater than 0
	curr    T
}

func newIntGen[T constraints.Integer](start, stop T, step uint64) *intGen[T] {
	if start == stop {
		panic("goloop: intGen: use intGenOne instead if start is equal to stop")
	}
	if step <= 0 {
		panic("goloop: intGen: step must be greater than 0")
	}
	var (
		kind         = reflect.TypeOf(start).Kind()
		panicStepErr = func(max uint64) {
			if step > max {
				panic(fmt.Errorf("goloop: intGen: step(%d) exceeds the maximum %s value", step, kind))
			}
		}
	)
	switch kind {
	case reflect.Int:
		panicStepErr(math.MaxInt)
	case reflect.Int8:
		panicStepErr(math.MaxInt8)
	case reflect.Int16:
		panicStepErr(math.MaxInt16)
	case reflect.Int32:
		panicStepErr(math.MaxInt32)
	case reflect.Int64:
		panicStepErr(math.MaxInt64)
	case reflect.Uint:
		panicStepErr(math.MaxUint)
	case reflect.Uint8:
		panicStepErr(math.MaxUint8)
	case reflect.Uint16:
		panicStepErr(math.MaxUint16)
	case reflect.Uint32:
		panicStepErr(math.MaxUint32)
	case reflect.Uint64:
		panicStepErr(math.MaxUint64)
	default:
		panic(fmt.Errorf("goloop: intGen: unsupported types (%s)", kind))
	}
	return &intGen[T]{
		start: start,
		stop:  stop,
		step:  T(step),
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
		if g.curr+g.step < g.stop {
			g.curr += g.step
			return true
		}
		return false
	} else {
		if g.curr-g.step > g.curr { // check overflow
			return false
		}
		if g.curr-g.step > g.stop {
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

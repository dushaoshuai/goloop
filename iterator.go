package goloop

import (
	"golang.org/x/exp/constraints"
)

type Integer interface {
	constraints.Integer | byte
}

// I is the conventional iteration variable i.
type I[T Integer] struct {
	// I is the value of the iteration variable i.
	I T
	// Break breaks the loop.
	Break func()
}

type iterator[T Integer] struct {
	// c is used to communicate iteration values.
	c chan I[T]
	// Close breakChan to signal that it's time to break the loop.
	breakChan chan struct{}
}

func newIterator[T Integer]() *iterator[T] {
	return &iterator[T]{
		c:         make(chan I[T]),
		breakChan: make(chan struct{}),
	}
}

func (i *iterator[T]) breakFunc() {
	// Allow breakFunc to be called multiple times.
	defer func() {
		recover()
	}()

	close(i.breakChan)

	// Empty channel i.c to not immediately return to the caller of breakFunc
	// until i.c is closed. Prevent the caller from receiving more iteration
	// values after calling breakFunc.
	for range i.c {
	}
}

func (i *iterator[T]) finish() {
	close(i.c)
}

func (i *iterator[T]) iter(value T) (breaked bool) {
	select {
	case <-i.breakChan:
		return true
	case i.c <- I[T]{I: value, Break: i.breakFunc}:
		return false
	}
}

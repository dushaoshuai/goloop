package goloop

// I is the conventional iteration variable i.
type I[T constraint] struct {
	// I is the value of the iteration variable i.
	I T
	// Break breaks the loop.
	Break func()
}

type chanIter[T constraint] struct {
	// c is used to communicate iteration values.
	c chan I[T]
	// Close breakChan to signal that it's time to break the loop.
	breakChan chan struct{}
}

func newIterator[T constraint]() *chanIter[T] {
	return &chanIter[T]{
		c:         make(chan I[T]),
		breakChan: make(chan struct{}),
	}
}

func (i *chanIter[T]) breakFunc() {
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

func (i *chanIter[T]) finish() {
	close(i.c)
}

func (i *chanIter[T]) iter(value T) (breaked bool) {
	select {
	case <-i.breakChan:
		return true
	case i.c <- I[T]{I: value, Break: i.breakFunc}:
		return false
	}
}

type sliceIter[T constraint] []T

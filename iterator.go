package goloop

// I is the conventional iteration variable i.
type I struct {
	// I is the value of the iteration variable i.
	I int
	// Break breaks the loop.
	// For now Break has no difference between each iteration.
	Break func()
}

type iterator struct {
	// c is used to communicate iteration values.
	c chan I
	// close breakChan to signal that it's time to break the loop.
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

	// Empty channel i.c to not immediately return to the caller of breakFunc
	// until i.c is closed. Prevent the caller from receiving more iteration
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

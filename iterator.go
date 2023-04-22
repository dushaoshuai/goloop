package goloop

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

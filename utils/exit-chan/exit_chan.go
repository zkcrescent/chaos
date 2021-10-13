package exitChan

import (
	"sync/atomic"
)

type ExitChan struct {
	exit   chan struct{}
	exited int32
}

func New() *ExitChan {
	return &ExitChan{
		exit: make(chan struct{}),
	}
}

func (ec *ExitChan) Chan() <-chan struct{} {
	return ec.exit
}

func (ec *ExitChan) Close() {
	if !atomic.CompareAndSwapInt32(&ec.exited, 0, 1) {
		return
	}
	close(ec.exit)
}

func (ec *ExitChan) Exited() bool {
	return atomic.LoadInt32(&ec.exited) == 1
}

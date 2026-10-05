package once

import (
	"primitives/internal/futex"
	"sync/atomic"
)

type Once struct {
	state uint32
}

const (
	notStarted = iota
	inProgress
	done
)

func (o *Once) Do(f func()) {
	for {
		state := atomic.LoadUint32(&o.state)
		switch state {
		case done:
			return
		case inProgress:
			futex.Wait(&o.state, inProgress)
			continue
		case notStarted:
			if atomic.CompareAndSwapUint32(&o.state, notStarted, inProgress) {
				defer func() {
					atomic.StoreUint32(&o.state, done)
					futex.WakeAll(&o.state)
				}()
				f()
				return
			}
			continue
		}
	}
}

func (o *Once) Done() bool {
	return atomic.LoadUint32(&o.state) == done
}

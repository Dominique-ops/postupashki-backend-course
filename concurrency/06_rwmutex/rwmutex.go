package rwmutex

import (
	"primitives/internal/futex"
	"sync/atomic"
)

const writer = 1 << 31

type RWMutex struct {
	state          uint32
	writer_waiting uint32
}

func (rw *RWMutex) RLock() {
	for {
		if atomic.LoadUint32(&rw.writer_waiting) > 0 {
			state := atomic.LoadUint32(&rw.state)
			futex.Wait(&rw.state, state)
			continue
		}

		state := atomic.LoadUint32(&rw.state)

		if state&writer != 0 {
			futex.Wait(&rw.state, state)
			continue
		}

		if atomic.CompareAndSwapUint32(&rw.state, state, state+1) {
			return
		}
	}
}

func (rw *RWMutex) RUnlock() {
	for {
		state := atomic.LoadUint32(&rw.state)
		if state == 0 || state == writer {
			panic("rwmutex: RUnlock of unlocked RWMutex")
		}
		if atomic.CompareAndSwapUint32(&rw.state, state, state-1) {
			if state == 1 {
				futex.WakeAll(&rw.state)
			}
			return
		}
	}
}

func (rw *RWMutex) Lock() {
	atomic.AddUint32(&rw.writer_waiting, 1)
	for {

		if !atomic.CompareAndSwapUint32(&rw.state, 0, writer) {
			futex.Wait(&rw.state, writer)
			continue
		}
		break

	}

	atomic.AddUint32(&rw.writer_waiting, ^uint32(0))
}

func (rw *RWMutex) Unlock() {

	if !atomic.CompareAndSwapUint32(&rw.state, writer, 0) {
		panic("rwmutex: Unlock of unlocked RWMutex")
	}
	futex.WakeAll(&rw.state)

}

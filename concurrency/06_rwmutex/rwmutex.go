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
		waiting := atomic.LoadUint32(&rw.writer_waiting)
		if waiting > 0 {
			// ВАЖНО: ждём изменения writer_waiting,
			// а не state.
			futex.Wait(&rw.writer_waiting, waiting)
			continue
		}

		state := atomic.LoadUint32(&rw.state)

		if state&writer != 0 {
			futex.Wait(&rw.state, state)
			continue
		}

		if atomic.LoadUint32(&rw.writer_waiting) > 0 {
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
		state := atomic.LoadUint32(&rw.state)

		if state == 0 {
			if atomic.CompareAndSwapUint32(
				&rw.state,
				0,
				writer,
			) {
				break
			}

			continue
		}

		futex.Wait(&rw.state, state)
	}

	waiting := atomic.AddUint32(
		&rw.writer_waiting,
		^uint32(0),
	)

	if waiting == 0 {
		futex.WakeAll(&rw.writer_waiting)
	}
}

func (rw *RWMutex) Unlock() {

	if !atomic.CompareAndSwapUint32(&rw.state, writer, 0) {
		panic("rwmutex: Unlock of unlocked RWMutex")
	}
	futex.WakeAll(&rw.state)

}

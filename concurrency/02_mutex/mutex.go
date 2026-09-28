package mutex

import (
	"primitives/internal/futex"
	"sync/atomic"
)

const (
	free = iota
	held
	contended
)

type Mutex struct {
	state uint32
}

func (m *Mutex) Lock() {
	if atomic.CompareAndSwapUint32(&m.state, free, held) {
		return
	}
	for {
		for i := 0; i < 100; i++ {
			if m.TryLock() {
				return
			}
		}
		atomic.CompareAndSwapUint32(&m.state, held, contended)
		futex.Wait(&m.state, contended)

		if atomic.SwapUint32(&m.state, contended) == free {
			return
		}
	}
}

func (m *Mutex) TryLock() bool {
	return atomic.CompareAndSwapUint32(&m.state, free, held)
}

func (m *Mutex) Unlock() {
	old := atomic.SwapUint32(&m.state, free)

	if old == free {
		panic("mutex is not locked")
	}

	if old == contended {
		futex.Wake(&m.state)
	}
}

package semaphore

import (
	"primitives/internal/futex"
	"sync/atomic"
)

const (
	free = iota
	held
	contended
)

type Semaphore struct {
	permits uint32
	waiters uint32
}

func New(n int) *Semaphore {
	if n < 0 {
		panic("semaphore: negative permits")
	}

	if uint64(n) > uint64(^uint32(0)) {
		panic("semaphore: too many permits")
	}

	return &Semaphore{
		permits: uint32(n),
	}
}
func (s *Semaphore) Acquire() {
	for {
		old := atomic.LoadUint32(&s.permits)

		if old > 0 {
			if atomic.CompareAndSwapUint32(&s.permits, old, old-1) {
				return
			}
			continue
		}

		break
	}

	atomic.AddUint32(&s.waiters, 1)

	for {
		old := atomic.LoadUint32(&s.permits)

		if old > 0 {
			if atomic.CompareAndSwapUint32(&s.permits, old, old-1) {
				atomic.AddUint32(&s.waiters, ^uint32(0))
				return
			}
			continue
		}

		futex.Wait(&s.permits, 0)
	}
}

func (s *Semaphore) TryAcquire() bool {
	for {
		old := atomic.LoadUint32(&s.permits)
		if old == 0 {
			return false
		}
		if atomic.CompareAndSwapUint32(&s.permits, old, old-1) {
			return true
		}
	}
}

func (s *Semaphore) Release() {
	atomic.AddUint32(&s.permits, 1)

	if atomic.LoadUint32(&s.waiters) > 0 {
		futex.Wake(&s.permits)
	}
}

func (s *Semaphore) Available() int {
	return int(atomic.LoadUint32(&s.permits))
}

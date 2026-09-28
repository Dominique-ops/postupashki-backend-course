package semaphore

import (
	"primitives/internal/futex"
	"sync/atomic"
)

type Semaphore struct {
	permits uint32
}

func New(n int) *Semaphore {
	return &Semaphore{
		permits: uint32(n),
	}
}

func (s *Semaphore) Acquire() {
	for {

		if old := atomic.LoadUint32(&s.permits); old > 0 {
			if atomic.CompareAndSwapUint32(&s.permits, old, old-1) {
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
	futex.Wake(&s.permits)
}

func (s *Semaphore) Available() int {
	return int(atomic.LoadUint32(&s.permits))
}

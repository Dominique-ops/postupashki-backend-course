package waitgroup

import (
	"sync/atomic"

	"primitives/internal/futex"
)

type WaitGroup struct {
	count uint32
}

func (wg *WaitGroup) Add(delta int) {
	for {
		old := atomic.LoadUint32(&wg.count)

		newCount := int64(old) + int64(delta)

		if newCount < 0 {
			panic("negative WaitGroup counter")
		}

		if newCount > int64(^uint32(0)) {
			panic("WaitGroup counter overflow")
		}

		if atomic.CompareAndSwapUint32(
			&wg.count,
			old,
			uint32(newCount),
		) {
			if newCount == 0 && old != 0 {
				futex.WakeAll(&wg.count)
			}
			return
		}
	}
}

func (wg *WaitGroup) Done() {
	wg.Add(-1)
}

func (wg *WaitGroup) Wait() {
	for {
		n := atomic.LoadUint32(&wg.count)

		if n == 0 {
			return
		}

		futex.Wait(&wg.count, n)
	}
}

package waitgroup

import (
	"primitives/internal/futex"
	"sync/atomic"
)

type WaitGroup struct {
	count uint32
}

func (wg *WaitGroup) Add(delta int) {
	atomic.AddUint32(&wg.count, uint32(delta))
}

func (wg *WaitGroup) Done() {
	for {
		old := atomic.LoadUint32(&wg.count)

		if old == 0 {
			panic("negative WaitGroup counter")
		}

		if atomic.CompareAndSwapUint32(&wg.count, old, old-1) {
			if old == 1 {
				futex.WakeAll(&wg.count)
			}
			return
		}
	}
}
func (wg *WaitGroup) Wait() {
	for {

		if n := atomic.LoadUint32(&wg.count); n != 0 {
			futex.Wait(&wg.count, n)
		} else {
			return
		}

	}
}

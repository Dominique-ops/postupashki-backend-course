package barrier

import (
	"primitives/internal/futex"
	"sync/atomic"
)

type Barrier struct {
	need    uint32
	arrived uint32
	round   uint32
}

func New(n int) *Barrier {
	return &Barrier{
		need: uint32(n),
	}
}

func (b *Barrier) Wait() {
	myRound := atomic.LoadUint32(&b.round)

	arrived := atomic.AddUint32(&b.arrived, 1)

	if arrived == b.need {
		atomic.StoreUint32(&b.arrived, 0)
		atomic.AddUint32(&b.round, 1)

		futex.WakeAll(&b.round)
		return
	}

	for atomic.LoadUint32(&b.round) == myRound {
		futex.Wait(&b.round, myRound)
	}
}

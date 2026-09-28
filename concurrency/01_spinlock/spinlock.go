package spinlock

import "sync/atomic"

type Spinlock struct {
	locked atomic.Bool
}

func (s *Spinlock) Lock() {
	for s.locked.Swap(true) {
	}
}

func (s *Spinlock) TryLock() bool {
	return !s.locked.Swap(true)
}

func (s *Spinlock) Unlock() {
	if !s.locked.Swap(false) {
		panic("unlock of unlocked spinlock")
	}
}

type TTAS struct {
	locked atomic.Bool
}

func (s *TTAS) Lock() {
	for {
		for s.locked.Load() {

		}
		if !s.locked.Swap(true) {
			return
		}
	}

}

func (s *TTAS) TryLock() bool {
	return !s.locked.Swap(true)
}

func (s *TTAS) Unlock() {
	if !s.locked.Swap(false) {
		panic("unlock of unlocked spinlock")
	}
}

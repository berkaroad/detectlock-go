package detectlock

import (
	"reflect"
	"sync"
)

// Mutex wrapped from sync.Mutex
type Mutex struct {
	m sync.Mutex
}

func (l *Mutex) Lock() {
	if debug {
		lockerPtr := reflect.ValueOf(l).Pointer()
		acquire(lockerPtr, false, l.m.Lock)
	} else {
		l.m.Lock()
	}
}

func (l *Mutex) TryLock() bool {
	if debug {
		lockerPtr := reflect.ValueOf(l).Pointer()
		return tryAcquire(lockerPtr, false, l.m.TryLock)
	} else {
		return l.m.TryLock()
	}
}

func (l *Mutex) Unlock() {
	if debug {
		lockerPtr := reflect.ValueOf(l).Pointer()
		release(lockerPtr, false, l.m.Unlock)
	} else {
		l.m.Unlock()
	}
}

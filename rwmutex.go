package detectlock

import (
	"reflect"
	"sync"
)

// Mutex wrapped from sync.RWMutex
type RWMutex struct {
	rwm sync.RWMutex
}

func (l *RWMutex) RLock() {
	if debug {
		lockerPtr := reflect.ValueOf(l).Pointer()
		acquire(lockerPtr, true, l.rwm.RLock)
	} else {
		l.rwm.RLock()
	}
}

func (l *RWMutex) TryRLock() bool {
	if debug {
		lockerPtr := reflect.ValueOf(l).Pointer()
		return tryAcquire(lockerPtr, false, l.rwm.TryRLock)
	} else {
		return l.rwm.TryRLock()
	}
}

func (l *RWMutex) RUnlock() {
	if debug {
		lockerPtr := reflect.ValueOf(l).Pointer()
		release(lockerPtr, true, l.rwm.RUnlock)
	} else {
		l.rwm.RUnlock()
	}
}

func (l *RWMutex) Lock() {
	if debug {
		lockerPtr := reflect.ValueOf(l).Pointer()
		acquire(lockerPtr, false, l.rwm.Lock)
	} else {
		l.rwm.Lock()
	}
}

func (l *RWMutex) TryLock() bool {
	if debug {
		lockerPtr := reflect.ValueOf(l).Pointer()
		return tryAcquire(lockerPtr, false, l.rwm.TryLock)
	} else {
		return l.rwm.TryLock()
	}
}

func (l *RWMutex) Unlock() {
	if debug {
		lockerPtr := reflect.ValueOf(l).Pointer()
		release(lockerPtr, false, l.rwm.Unlock)
	} else {
		l.rwm.Unlock()
	}
}

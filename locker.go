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

func (l *Mutex) Unlock() {
	if debug {
		lockerPtr := reflect.ValueOf(l).Pointer()
		release(lockerPtr, false, l.m.Unlock)
	} else {
		l.m.Unlock()
	}
}

// Mutex wrapped from sync.RWMutex
type RWMutex struct {
	rwm sync.RWMutex
}

func (l *RWMutex) Lock() {
	if debug {
		lockerPtr := reflect.ValueOf(l).Pointer()
		acquire(lockerPtr, false, l.rwm.Lock)
	} else {
		l.rwm.Lock()
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

func (l *RWMutex) RLock() {
	if debug {
		lockerPtr := reflect.ValueOf(l).Pointer()
		acquire(lockerPtr, true, l.rwm.RLock)
	} else {
		l.rwm.RLock()
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

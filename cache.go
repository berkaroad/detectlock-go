package detectlock

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
)

const shardCount uint64 = 1 << 16

// Status of locker.
const (
	StatusWait     byte = iota // wait or r-wait
	StatusAcquired             // acquired or r-acquired
)

var cache cacheMap

// LockerState the state of locker.
type LockerState struct {
	LockerPtr uintptr
	Status    byte
	RLocker   bool
	Caller    *runtime.Frame
}

// String of LockerState, format: (<locker-id>, <locker-status>)
func (l LockerState) String() string {
	status := "wait"
	if l.Status == StatusAcquired {
		status = "acquired"
	}
	if l.RLocker {
		status = "r-" + status
	}

	stackInfo := ""
	if l.Caller != nil {
		stackInfo = fmt.Sprintf("%s (file: %s:%d)", l.Caller.Function, l.Caller.File, l.Caller.Line)
	}
	return fmt.Sprintf("(%#x, %s, %s)", l.LockerPtr, status, stackInfo)
}

// LockerStateList the list of LockerState
type LockerStateList []LockerState

func (l LockerStateList) Len() int {
	return len(l)
}

func (l LockerStateList) Less(i, j int) bool {
	return l[i].Status > l[j].Status
}

func (l LockerStateList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// String of LockerStateList, format: [(<locker-id>, <locker-status>), ...]
func (l LockerStateList) String() string {
	if len(l) == 0 {
		return "[]"
	}
	sb := &strings.Builder{}
	sb.WriteString("[")
	llen := len(l)
	for i, locker := range l {
		sb.WriteString(locker.String())
		if i < llen-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

// Items of goroutine that use locker.
func Items() map[uint64]LockerStateList {
	items := make(map[uint64]LockerStateList)
	for _, mapShard := range cache {
		func() {
			defer mapShard.locker.Unlock()
			mapShard.locker.Lock()
			for k, v := range mapShard.items {
				lockers := make(LockerStateList, len(v))
				for i := 0; i < len(v); i++ {
					lockers[i] = *v[i]
				}
				items[k] = lockers
			}
		}()
	}
	return items
}

type cacheMap []*cacheMapShard

type cacheMapShard struct {
	locker sync.RWMutex
	items  map[uint64][]*LockerState
}

func clear() {
	cache = make(cacheMap, shardCount)
}

func reset() {
	cache = make(cacheMap, shardCount)
	var i uint64 = 0
	for ; i < shardCount; i++ {
		cache[i] = &cacheMapShard{items: make(map[uint64][]*LockerState)}
	}
}

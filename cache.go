package detectlock

import (
	"fmt"
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

func init() {
	cache = make(cacheMap, shardCount)
	var i uint64 = 0
	for ; i < shardCount; i++ {
		cache[i] = &cacheMapShard{items: make(map[uint64][]*LockerState)}
	}
}

// LockerState the state of locker.
type LockerState struct {
	LockerPtr uintptr
	Status    byte
	RLocker   bool
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
	return fmt.Sprintf("(%#x, %s)", l.LockerPtr, status)
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

func acquire(lockerPtr uintptr, rLocker bool, doLock func()) {
	gid := getGoroutineID()
	shardKey := gid % shardCount
	mapShard := cache[shardKey]
	var locker *LockerState

	func() {
		defer mapShard.locker.Unlock()
		mapShard.locker.Lock()
		var lockers []*LockerState
		if exists, ok := mapShard.items[gid]; ok {
			lockers = exists
		} else {
			lockers = make([]*LockerState, 0)
		}
		locker = &LockerState{LockerPtr: lockerPtr, Status: StatusWait, RLocker: rLocker}
		lockers = append(lockers, locker)
		mapShard.items[gid] = lockers
	}()
	doLock()
	locker.Status = StatusAcquired
}

func release(lockerPtr uintptr, rLocker bool, doUnlock func()) {
	doUnlock()

	gid := getGoroutineID()
	shardKey := gid % shardCount
	mapShard := cache[shardKey]

	defer mapShard.locker.Unlock()
	mapShard.locker.Lock()
	if lockers, ok := mapShard.items[gid]; ok {
		removeIndex := -1
		for i := 0; i < len(lockers); i++ {
			l := lockers[i]
			if l.LockerPtr == lockerPtr && l.Status == StatusAcquired && l.RLocker == rLocker {
				removeIndex = i
				break
			}
		}
		if removeIndex < 0 {
			return
		}
		llen := len(lockers)
		if llen == 1 {
			lockers = nil
		} else {
			lockers = append(lockers[:removeIndex], lockers[removeIndex+1:]...)
		}
		if len(lockers) == 0 {
			delete(mapShard.items, gid)
		} else {
			mapShard.items[gid] = lockers
		}
	}
}

func clear() {
	var i uint64 = 0
	for ; i < shardCount; i++ {
		cache[i].items = make(map[uint64][]*LockerState)
	}
}

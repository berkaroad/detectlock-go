package detectlock

import (
	"fmt"
	"sort"
	"strings"
)

var debug bool

// EnableDebug enable for debug
func EnableDebug() {
	reset()
	debug = true
}

// DisableDebug disable for debug
func DisableDebug() {
	debug = false
	clear()
}

// DetectAcquired detect goroutine that locker acquired.
func DetectAcquired(items map[uint64]LockerStateList) GoroutineLockerList {
	l := make(GoroutineLockerList, 0)
	for gid, lockers := range items {
		lockerAcquired := false
		for _, locker := range lockers {
			if locker.Status == StatusAcquired {
				lockerAcquired = true
			}
		}
		if lockerAcquired {
			l = append(l, GoroutineLocker{GoroutineID: gid, Lockers: lockers})
		}
	}
	sort.Sort(l)
	return l
}

// DetectLockedEachOther detect goroutine that locked each other.
func DetectLockedEachOther(items map[uint64]LockerStateList) GoroutineLockerList {
	l := make(GoroutineLockerList, 0)
	for gid, lockers := range items {
		sortedLockers := make(LockerStateList, len(lockers))
		copy(sortedLockers, lockers)
		lockerAcquired := false
		acquiredLockers := make([]uintptr, 0)
		waiting := false
		for _, locker := range sortedLockers {
			if locker.Status == StatusAcquired {
				acquiredLockers = append(acquiredLockers, locker.LockerPtr)
				if !locker.RLocker {
					lockerAcquired = true
				} else {
				loopOtherGID:
					for ogid, olockers := range items {
						if ogid == gid {
							continue
						}
						waitByOtherGID := false
						lockerAcquiredByOtherGID := false
						for _, olocker := range olockers {
							if olocker.LockerPtr == locker.LockerPtr && olocker.Status == StatusWait && !olocker.RLocker {
								waitByOtherGID = true
							}
							if olocker.Status == StatusAcquired {
								lockerAcquiredByOtherGID = true
							}
							if waitByOtherGID && lockerAcquiredByOtherGID {
								lockerAcquired = true
								break loopOtherGID
							}
						}
					}
				}
			} else {
				existsInAcquiredLockers := false
				for _, acuqiredLocker := range acquiredLockers {
					if acuqiredLocker == locker.LockerPtr {
						existsInAcquiredLockers = true
						break
					}
				}
				if !existsInAcquiredLockers {
					waiting = true
				}
			}
		}
		if lockerAcquired && waiting {
			l = append(l, GoroutineLocker{GoroutineID: gid, Lockers: lockers})
		}
	}
	sort.Sort(l)
	return l
}

// DetectReentry detect goroutine that reentry locker occurred.
func DetectReentry(items map[uint64]LockerStateList) GoroutineLockerList {
	l := make(GoroutineLockerList, 0)
	for gid, lockers := range items {
		sortedLockers := make(LockerStateList, len(lockers))
		copy(sortedLockers, lockers)
		acquiredLockers := make([]uintptr, 0)
		for _, locker := range sortedLockers {
			if locker.Status == StatusAcquired {
				acquiredLockers = append(acquiredLockers, locker.LockerPtr)
			} else {
				existsInAcquiredLockers := false
				for _, acuqiredLocker := range acquiredLockers {
					if acuqiredLocker == locker.LockerPtr {
						existsInAcquiredLockers = true
						break
					}
				}
				if existsInAcquiredLockers {
					l = append(l, GoroutineLocker{GoroutineID: gid, Lockers: lockers})
				}
			}
		}
	}
	sort.Sort(l)
	return l
}

// GoroutineLocker goroutine with lockers.
type GoroutineLocker struct {
	GoroutineID uint64
	Lockers     LockerStateList
}

// String of GoroutineLocker, format: goroutine <gid>: [(<locker-id>, <locker-status>), ...]\n, like:
//
// goroutine 53: [(0xc000014080, acquired), (0xc000014088, wait)]
func (l GoroutineLocker) String() string {
	return fmt.Sprintf("goroutine %d: %s\n", l.GoroutineID, l.Lockers)
}

// GoroutineLockerList the list of GoroutineLocker
type GoroutineLockerList []GoroutineLocker

func (l GoroutineLockerList) Len() int {
	return len(l)
}

func (l GoroutineLockerList) Less(i, j int) bool {
	return l[i].GoroutineID < l[j].GoroutineID
}

func (l GoroutineLockerList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// String of GoroutineLockerList, format: goroutine <gid>: [(<locker-id>, <locker-status>), ...]\n..., like:
//
// goroutine 53: [(0xc000014080, acquired), (0xc000014088, wait)]
//
// goroutine 54: [(0xc000014088, acquired), (0xc000014080, wait)]
func (l GoroutineLockerList) String() string {
	if len(l) == 0 {
		return ""
	}
	sb := &strings.Builder{}
	for _, glocker := range l {
		sb.WriteString(glocker.String())
	}
	return sb.String()
}

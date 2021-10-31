// Package detectlock for detect dead locker
//
// 1. replace mutex locker:
//
// 1) replace "var locker *sync.Mutex = &sync.Mutex{}" to "var locker *detectlock.Mutex = &detectlock.Mutex{}"
//
// 2) replace "var locker *sync.RWMutex = &sync.RWMutex{}" to "var locker *detectlock.RWMutex = &detectlock.RWMutex{}"
//
// 2. enable debug on startup
//
// detectlock.EnableDebug()" or disable it by "detectlock.DisableDebug()"
//
// 3. detect dead locker
//
// items := detectlock.Items()
//
// detect reentry locker: "detectlock.DetectReentry(items)"
//
// detect locked each other: "detectlock.DetectLockedEachOther(items)"
//
// detect acquired owners: "detectlock.DetectAcquired(items)"
package detectlock

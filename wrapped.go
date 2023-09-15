package detectlock

func acquire(lockerPtr uintptr, rLocker bool, doLock func()) {
	if doLock == nil {
		return
	}

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

	locker.Caller = getCaller(4)
	doLock()
	locker.Status = StatusAcquired
}

func tryAcquire(lockerPtr uintptr, rLocker bool, tryLock func() bool) bool {
	if tryLock == nil {
		return false
	}

	if tryLock() {
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
			locker = &LockerState{LockerPtr: lockerPtr, Status: StatusAcquired, RLocker: rLocker}
			lockers = append(lockers, locker)
			mapShard.items[gid] = lockers
		}()

		locker.Caller = getCaller(4)
		return true
	} else {
		return false
	}
}

func release(lockerPtr uintptr, rLocker bool, doUnlock func()) {
	if doUnlock == nil {
		return
	}
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

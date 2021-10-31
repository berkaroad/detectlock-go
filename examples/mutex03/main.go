package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/berkaroad/detectlock-go/detectlock"
)

var Locker4 *detectlock.RWMutex
var Locker5 sync.Locker

func init() {
	Locker4 = &detectlock.RWMutex{}
	Locker5 = &detectlock.Mutex{}
}

func D() {
	defer Locker4.RUnlock()
	Locker4.RLock()
	time.Sleep(time.Millisecond * 100)
	defer Locker5.Unlock()
	Locker5.Lock()
	time.Sleep(time.Millisecond * 400)
}
func E() {
	time.Sleep(time.Millisecond * 100)
	defer Locker5.Unlock()
	Locker5.Lock()
	time.Sleep(time.Millisecond * 100)
	defer Locker4.Unlock()
	Locker4.Lock()
	time.Sleep(time.Millisecond * 400)
}

func main() {
	detectlock.EnableDebug()
	for i := 0; i < 10; i++ {
		go D()
		go E()
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		time.Sleep(time.Second * 2)

		items := detectlock.Items()

		fmt.Println("--- DetectAcquired ---")
		fmt.Println(detectlock.DetectAcquired(items))

		fmt.Println("--- DetectLockedEachOther ---")
		fmt.Println(detectlock.DetectLockedEachOther(items))

		// fmt.Println("--- all stack ---")
		// b := make([]byte, 102400)
		// b = b[:runtime.Stack(b, true)]
		//fmt.Println(string(b))
		// fmt.Println("------- end --------")

		detectlock.DisableDebug()
		wg.Done()
	}()
	wg.Wait()
}

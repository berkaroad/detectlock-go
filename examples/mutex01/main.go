package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/berkaroad/detectlock-go/detectlock"
)

var Locker1 sync.Locker
var Locker2 sync.Locker

func init() {
	Locker1 = &detectlock.Mutex{}
	Locker2 = &detectlock.Mutex{}
}
func A() {
	defer Locker1.Unlock()
	Locker1.Lock()
	time.Sleep(time.Millisecond * 300)
	defer Locker2.Unlock()
	Locker2.Lock()

	time.Sleep(time.Millisecond * 300)
}

func B() {
	defer Locker2.Unlock()
	Locker2.Lock()
	time.Sleep(time.Millisecond * 400)
	defer Locker1.Unlock()
	Locker1.Lock()

	time.Sleep(time.Millisecond * 400)
}

func main() {
	detectlock.EnableDebug()
	for i := 0; i < 10; i++ {
		go A()
		go B()
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

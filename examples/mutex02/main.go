package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/berkaroad/detectlock-go/detectlock"
)

var Locker3 sync.Locker

func init() {
	Locker3 = &detectlock.Mutex{}
}

func C() {
	defer Locker3.Unlock()
	Locker3.Lock()

	time.Sleep(time.Millisecond * 400)
	Locker3.Lock()
}

func main() {
	detectlock.EnableDebug()
	for i := 0; i < 10; i++ {
		go C()
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		time.Sleep(time.Second * 2)

		items := detectlock.Items()

		fmt.Println("--- DetectAcquired ---")
		fmt.Println(detectlock.DetectAcquired(items))

		fmt.Println("--- DetectReentry ---")
		fmt.Println(detectlock.DetectReentry(items))

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

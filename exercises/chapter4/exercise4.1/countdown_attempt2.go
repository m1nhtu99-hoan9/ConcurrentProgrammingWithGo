package main

import (
	"fmt"
	"sync"
	"time"
)

// FALSE ATTEMPT!! There is only decreasing `count` without increasing -> wait without signal -> meaningless

func countdown2(seconds *int, cond *sync.Cond) {
	cond.L.Lock()
	if *seconds <= 1 {
		cond.Wait()
	}

	cond.Signal()
	time.Sleep(time.Second)
	*seconds -= 1
	cond.L.Unlock()
}

func main() {
	count := 5
	cond := sync.NewCond(&sync.Mutex{})
	go countdown2(&count, cond)

	for count > 0 {
		time.Sleep(500 * time.Millisecond)
		fmt.Println(count)
	}
	cond.L.Unlock()
}

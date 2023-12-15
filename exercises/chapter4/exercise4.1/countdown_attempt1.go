package main

import (
	"fmt"
	"sync"
	"time"
)

/* IDEAS for improvements:
- `count` as semaphore
- `count` as atomic int
*/

func countdown(seconds *int, mutex *sync.RWMutex) {
	for {
		mutex.RLock()
		if *seconds <= 0 {
			mutex.RUnlock()
			break
		} else {
			mutex.RUnlock()
		}

		time.Sleep(1 * time.Second)
		mutex.Lock()
		*seconds -= 1
		mutex.Unlock()
	}
}

func main() {
	count := 5
	mutex := sync.RWMutex{}
	go countdown(&count, &mutex)
	for {
		mutex.RLock()

		if count > 0 {
			time.Sleep(500 * time.Millisecond)
			fmt.Println(count)
			mutex.RUnlock()
		} else {
			mutex.RUnlock()
			break
		}
	}
}

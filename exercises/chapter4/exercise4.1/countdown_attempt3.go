package main

import (
	"fmt"
	"sync"
	"time"
)

func countdownChan(countChan chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		val, ok := <-countChan
		if !ok || val <= 0 {
			close(countChan)
			return
		}

		time.Sleep(time.Second)
		val--
		countChan <- val
	}
}

// FAILED ATTEMPT: deadlocks
func main() {
	countChan := make(chan int, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go countdownChan(countChan, &wg)

	// Initialize the channel with the first value
	countChan <- 5

	for {
		select {
		case val := <-countChan:
			time.Sleep(500 * time.Millisecond)
			fmt.Println(val)
			if val <= 0 {
				close(countChan)
				wg.Wait() // Wait for the countdownChan goroutine to complete
				return
			}
		}
	}
}

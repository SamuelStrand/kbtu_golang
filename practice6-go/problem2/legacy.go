package main

import (
	"fmt"
	"sync"
)

// 1-sentence explanation:
// The final answer is not 1000 because `counter++` is a non-atomic read-modify-write operation,
// so concurrent goroutines race and some increments are lost.

func main() {
	var counter int
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++ // data race
		}()
	}

	wg.Wait()
	fmt.Println(counter)
}

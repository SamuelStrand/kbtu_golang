package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	unsafeMap := make(map[string]int)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			unsafeMap["key"] = v
		}(i)
	}

	time.Sleep(10 * time.Millisecond)

	value := unsafeMap["key"]
	fmt.Printf("Value: %d\n", value)

	wg.Wait()
}

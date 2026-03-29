package main

import (
	"fmt"
	"sync"
)

func main() {
	var m sync.Map
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			m.Store("key", v)
		}(i)
	}

	wg.Wait()

	valueAny, _ := m.Load("key")
	value, _ := valueAny.(int)
	fmt.Printf("Value: %d\n", value)
}

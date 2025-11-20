//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"sync"
)

// go run -race race_condition.go
func main() {
	var count int
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			count++
		}()
	}

	wg.Wait()
	fmt.Println("Count:", count)
}

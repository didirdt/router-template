//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"sync"
)

type Data struct {
	Balance int
	Rmutex  sync.RWMutex
}

func main() {
	loop := 1000
	ic := make(chan map[string]int, loop) // Buffered channel to prevent blocking
	var v Data
	var wg sync.WaitGroup

	// Launch goroutines
	for range make([]struct{}, loop) {
		wg.Add(1)
		go f(&v, &wg, ic)
	}

	// Wait for all goroutines to complete in separate goroutine
	go func() {
		wg.Wait()
		close(ic) // Close channel after all writers are done
	}()

	// Process results from channel
	count := 0
	for result := range ic {
		if (result["value"] % 500) == 0 {
			fmt.Println("__________")
			fmt.Println("Read : ", result)
		}
		count++
		if count == 1000 { // We know we expect 1000 results
			break
		}
	}

	// Get final balance safely
	v.Rmutex.RLock()
	finalBalance := v.Balance
	v.Rmutex.RUnlock()

	fmt.Println("Finished", finalBalance)
}

func f(v *Data, wg *sync.WaitGroup, ch chan map[string]int) {
	defer wg.Done() // Use defer to ensure Done is called even if panic occurs

	v.Rmutex.Lock()
	v.Balance++
	currentBalance := v.Balance // Read balance while holding the lock
	v.Rmutex.Unlock()

	if (currentBalance % 500) == 0 {
		fmt.Println("ini : ", currentBalance)
	}

	// Send without blocking (channel is buffered)
	ch <- map[string]int{"value": currentBalance}
}

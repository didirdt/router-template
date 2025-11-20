//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"sync"
)

type SafeCounter struct {
	ch    chan int
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *SafeCounter) CounterValue() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}
func (c *SafeCounter) ChannelValue() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return <-c.ch
}
func (c *SafeCounter) Print() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ChannelValue()
}

func main() {
	counter := SafeCounter{}
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	fmt.Println("Final count:", counter.CounterValue()) // Always 100
}

package main

import (
	"fmt"
	"sync"
)

func main() {
	// Data source
	stringChan := make(chan string)
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := range stringChan {
			fmt.Printf("Listener 1: %v \n", i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := range stringChan {
			fmt.Printf("Listener 2: %v \n", i)
		}
	}()

	for i := 0; i < 5; i++ {
		go func() {
			stringChan <- "ini. string loading dulu " + fmt.Sprint(i)
		}()
	}

	stringChan <- "ini string"
	stringChan <- "ini juga string"

	close(stringChan)
	wg.Wait()
}

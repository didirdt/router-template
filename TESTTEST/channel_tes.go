//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"sync"
)

func main() {
	// Data source
	stringChan := make(chan string)
	var wg sync.WaitGroup
	arr := []map[string]any{}

	wg.Add(1)

	// stringChan <- "ini string"
	// stringChan <- "ini juga string"

	// go func() {
	// 	defer wg.Done()
	// 	for i := range stringChan {
	// 		fmt.Printf("Listener 1: %v \n", i)
	// 	}
	// }()
	for i := 0; i < 5; i++ {
		go listener(&wg, stringChan, 2, arr)
	}

	// for i := 0; i < 5; i++ {
	// 	func() {
	// 		stringChan <- "ini. string loading dulu " + fmt.Sprint(i)
	// 	}()
	// }

	stringChan <- "ini string"
	stringChan <- "ini juga string"

	close(stringChan)
	wg.Wait()

	fmt.Println("Array :", arr)
}

func listener(wg *sync.WaitGroup, stringChan chan string, id int, ar []map[string]any) {

}

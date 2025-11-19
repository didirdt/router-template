package main

import (
	"fmt"
	"sync"
)

func f(v *int, wg *sync.WaitGroup, m *sync.Mutex) {
	m.Lock()
	defer wg.Done()
	defer m.Unlock()
	*v++
	if (*v % 10) == 0 {
		fmt.Println("ini : ", *v)
		// time.Sleep(3 * time.Second)
	}
}
func fa(ch chan int, v *int) {
	ch <- *v
	close(ch)
}

func main() {
	// ic := make(chan int)
	var wg sync.WaitGroup
	var m sync.Mutex
	var v int = 0

	for i := 0; i < 100; i++ {
		wg.Add(1)
		// go fa(ic, &v)
		go f(&v, &wg, &m)
	}
	wg.Wait()

	// for v := range ic {
	// 	fmt.Println("__________")
	// 	fmt.Println("Read : ", v)
	// }
	fmt.Println("Finished", v)
}

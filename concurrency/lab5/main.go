package main

import (
	"fmt"
	"sync"
)

func localMultiply(vector []float64, ch chan float64) {
	result := 1.0
	for _, coord := range vector {
		result *= coord
	}

	ch <- result
}

func main() {
	var wg sync.WaitGroup

	ch := make(chan float64)

	vector1 := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	vector2 := []float64{6.0, 7.0, 8.0, 9.0, 10.0}

	wg.Add(1)
	go func() {
		defer wg.Done()
		localMultiply(vector1, ch)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		localMultiply(vector2, ch)
	}()

	go func() {
		wg.Wait()
		close(ch)
	}()

	skalMult := 0.0
	for local := range ch {
		skalMult += local
	}

	fmt.Println(skalMult)
}

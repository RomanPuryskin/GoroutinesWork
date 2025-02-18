package main

import (
	"fmt"
	"sync"
	"time"
)

func Count(number int, threadNumber int) {
	for i := number; i >= 0; i-- {
		fmt.Printf("Номер потока %d , текущее число: %d\n", threadNumber, i)
		time.Sleep(time.Second)
	}
}

func main() {
	var wg sync.WaitGroup

	var Number int

	fmt.Scan(&Number)

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func() {
			Count(Number, i)
			wg.Done()
		}()
	}

	wg.Wait()

}

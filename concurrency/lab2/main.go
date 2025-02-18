package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
)

type Array []float64

// метод для генерации массива заданной длины
func (arr *Array) GetRandomArray(length int) {
	for i := 0; i < length; i++ {
		*arr = append(*arr, randomFloatInRange(-100.0, 100.0))
	}
}

// функция для генерации рандомного вещественного числа
func randomFloatInRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func main() {

	var wg sync.WaitGroup

	var (
		arrayLength    int
		countOfThreads int
	)

	fmt.Scan(&arrayLength, &countOfThreads)

	threads := []*Array{}

	// заполним массив (указателями на пользовательский тип Array)
	for i := 0; i < countOfThreads; i++ {
		tempArr := &Array{}
		threads = append(threads, tempArr)
	}

	for index, thread := range threads {
		wg.Add(1)
		go func() {
			thread.GetRandomArray(arrayLength)
			fmt.Print("Номер потока в массиве: ", index, " ", *thread, " \n")
			wg.Done()
		}()
	}

	wg.Wait()

}

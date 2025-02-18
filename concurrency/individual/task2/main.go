package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// функция для которой считаем интеграл
func funtion(x float64) float64 {
	return math.Pow(x, 3)
}

func calculateIntegralConsist(a, b float64, N int) float64 {
	// шаг метода
	h := (b - a) / float64(N)

	// начальное значение
	result := (funtion(a) + funtion(b)) / 2

	for i := 1; i <= N-1; i++ {
		result += funtion(float64(i) * h)
	}

	result *= h

	return result

}

func calculateIntegralConcurr(a, b float64, N int, workers int) float64 {

	ch := make(chan float64, workers)
	var wg sync.WaitGroup
	var mtx sync.Mutex
	// шаг метода
	h := (b - a) / float64(N)

	// начальное значение
	result := (funtion(a) + funtion(b)) / 2

	// количество трапеций для вычисления на 1 горутину
	trapsForGoroutine := N / workers

	// запустим workers горутин
	for i := 0; i < workers; i++ {
		// для каждой горутины вычислим диапазон значений для подсчета
		start := i*trapsForGoroutine + 1
		var end int
		if i == workers-1 {
			end = N
		} else {
			end = start + trapsForGoroutine
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			localResult := 0.0
			for i := start; i < end; i++ {
				localResult += funtion(float64(i) * h)
			}
			ch <- localResult
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for localResult := range ch {
		mtx.Lock()
		result += localResult
		mtx.Unlock()
	}

	result *= h

	return result

}

func main() {
	// пределы интегрирования
	a := 0.0
	b := 2.0

	// количество трапеций
	N := 100000

	// количество горутин для параллельной работы
	workers := 6

	// время выполнения последовательного решения
	start1 := time.Now()
	answer1 := calculateIntegralConsist(a, b, N)
	duration1 := time.Since(start1)

	// время выполнения параллельного решения
	start2 := time.Now()
	answer2 := calculateIntegralConcurr(a, b, N, workers)
	duration2 := time.Since(start2)

	fmt.Println("Время выполнения послед: ", duration1, " Время выполнения паралл: ", duration2)
	fmt.Println(answer1, answer2)

}

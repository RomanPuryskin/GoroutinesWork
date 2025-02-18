package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Queue struct {
	Arr []int
	mtx sync.Mutex
}

// конструктор
func NewQueue() *Queue {
	return &Queue{
		Arr: []int{},
	}
}

// добавление элемента
func (q *Queue) Add(el int) {
	q.mtx.Lock()
	q.Arr = append(q.Arr, el)
	q.mtx.Unlock()
}

// удаление элемента
func (q *Queue) Delete() error {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	if len(q.Arr) == 0 {
		return errors.New("Нет элементов для удаления")
	}
	q.Arr = q.Arr[1:]

	return nil
}

// найти максимум
func (q *Queue) FindMax() (int, error) {
	if len(q.Arr) == 0 {
		return 0, errors.New("В очереди нет элементов")
	}

	curMax := (q.Arr)[0]

	for _, number := range q.Arr {
		if number > curMax {
			curMax = number
		}
	}

	return curMax, nil
}

// найти минимум
func (q *Queue) FindMin() (int, error) {
	if len(q.Arr) == 0 {
		return 0, errors.New("В очереди нет элементов")
	}

	curMax := (q.Arr)[0]

	for _, number := range q.Arr {
		if number < curMax {
			curMax = number
		}
	}

	return curMax, nil
}

// вывод содержимого
func (q *Queue) Print() {
	fmt.Println(q.Arr)
}

// функция реализующая барьерную синхронизацию
// рассмотрим 3 потока, каждый их которых заполяет очередь с разным временным интервалом
// дождемся пока все 3 потока добавят элемент в очередь чтобы перейти дальше
func BarrierAdding() {
	var wg sync.WaitGroup

	countThreads := 3 // количество потоков

	wg.Add(countThreads)

	q := NewQueue()

	for i := 1; i <= 5; i++ {
		go func() {
			// 1 поток (добавляет элемент за 8 секунд)
			defer wg.Done()
			time.Sleep(8 * time.Second)
			AddAndPrintInfo(q, 1, i)
		}()

		go func() {
			// 2 поток
			defer wg.Done()
			AddAndPrintInfo(q, 2, i)
		}()

		go func() {
			// 3 поток
			defer wg.Done()
			AddAndPrintInfo(q, 3, i)
		}()

		wg.Wait()
		fmt.Println("Все элементы добавлены")
		wg.Add(countThreads)

	}

}

func AddAndPrintInfo(q *Queue, threadNum int, elToAdd int) {
	q.Add(elToAdd)
	fmt.Println("В очередь добавлено число потоком ", threadNum, q.Arr)
}

func main() {
	BarrierAdding()
}

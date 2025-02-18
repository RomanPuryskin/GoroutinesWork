package main

import (
	"errors"
	"fmt"
	"sync"
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

func main() {
	var wg sync.WaitGroup

	queue := NewQueue()

	// добавляем элементы
	for i := 0; i <= 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			queue.Add(i)
		}()
	}

	// убираем элементы
	for i := 0; i <= 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			queue.Delete()
		}()
	}

	wg.Wait()

	queue.Print()
}

package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
	"time"
	"unicode"
)

// стрктура прочитанной информации

type ReadedInfo struct {
	FreqLetters  map[byte]int // для вычисления частот букв
	TotalLetters int          // для вычисления вероятности появления буквы
	mtx          sync.Mutex
}

// получить 5 самых распространенных букв
func (ri *ReadedInfo) Top5Letters() []string {
	ri.mtx.Lock()
	defer ri.mtx.Unlock()

	result := []string{}

	type structToSort struct {
		freq   int
		letter byte
	}

	arrToSort := []*structToSort{}
	for letter, freq := range ri.FreqLetters {
		arrToSort = append(arrToSort, &structToSort{freq, letter})
	}

	sort.SliceStable(arrToSort, func(i, j int) bool {
		if arrToSort[i].freq == arrToSort[j].freq {
			return arrToSort[i].letter > arrToSort[j].letter
		}
		return arrToSort[i].freq > arrToSort[j].freq
	})

	for i := 1; i <= 5; i++ {
		if len(ri.FreqLetters) >= i {
			result = append(result, string(arrToSort[i-1].letter))
		}
	}
	return result
}

// получить 3 самых редких буквы
func (ri *ReadedInfo) Last3Letters() []string {
	ri.mtx.Lock()
	defer ri.mtx.Unlock()
	result := []string{}

	type structToSort struct {
		freq   int
		letter byte
	}

	arrToSort := []*structToSort{}
	for letter, freq := range ri.FreqLetters {
		arrToSort = append(arrToSort, &structToSort{freq, letter})
	}

	sort.SliceStable(arrToSort, func(i, j int) bool {
		if arrToSort[i].freq == arrToSort[j].freq {
			return arrToSort[i].letter < arrToSort[j].letter
		}
		return arrToSort[i].freq < arrToSort[j].freq
	})

	for i := 1; i <= 3; i++ {
		if len(ri.FreqLetters) >= i {
			result = append(result, string(arrToSort[i-1].letter))
		}
	}
	return result
}

// вычисление вероятности появления буквы
func (ri *ReadedInfo) GetFreq(letter string) float64 {
	ri.mtx.Lock()
	defer ri.mtx.Unlock()
	return float64(ri.FreqLetters[[]byte(letter)[0]]) / float64(ri.TotalLetters)
}

// функция для последовательного вычисления
func ConsistentImplementation(ri *ReadedInfo, files []string) {
	for _, fileName := range files {
		file, _ := os.Open(fileName)
		defer file.Close()
		text, _ := io.ReadAll(file)
		for _, letter := range text {
			if unicode.IsLetter(rune(letter)) {
				ri.FreqLetters[letter]++
				ri.TotalLetters++
			}
		}
	}
}

// функция для параллельной реализации  Map Reduce
func ConcurrentImplementation(ri *ReadedInfo, files []string) {
	var wg sync.WaitGroup

	wg.Add(len(files))
	for i := 0; i < len(files); i++ {
		go func() {
			defer wg.Done()
			file, _ := os.Open(files[i])
			defer file.Close()
			text, _ := io.ReadAll(file)

			// создадим мапу для каждого потока
			localMap := make(map[byte]int)

			// заполним каждую локальную мапу своим текстом (Map)
			for _, letter := range text {
				if unicode.IsLetter(rune(letter)) {
					localMap[letter]++
				}
			}

			// сливаем каждую мапу параллельно в общую (Reduce)
			for key, value := range localMap {
				ri.mtx.Lock()
				ri.FreqLetters[key] += value
				ri.mtx.Unlock()

				ri.mtx.Lock()
				ri.TotalLetters += value
				ri.mtx.Unlock()
			}
		}()
	}

	wg.Wait()
}

func main() {

	files := []string{"file1.txt", "file2.txt", "file3.txt", "file4.txt", "file5.txt", "file6.txt", "file7.txt", "file8.txt", "file9.txt", "file10.txt"}

	ri1 := ReadedInfo{
		FreqLetters: make(map[byte]int),
	}

	ri2 := ReadedInfo{
		FreqLetters: make(map[byte]int),
	}

	// замерим время на последовательную реализацию
	start1 := time.Now()
	ConsistentImplementation(&ri1, files)
	duration1 := time.Since(start1)

	fmt.Println(ri1.Top5Letters())
	fmt.Println(ri1.Last3Letters())

	// замерим время на параллельную реализацию
	start2 := time.Now()
	ConcurrentImplementation(&ri2, files)
	duration2 := time.Since(start2)

	fmt.Println(ri2.Top5Letters())
	fmt.Println(ri2.Last3Letters())
	fmt.Println("Послед: ", duration1, " Паралл: ", duration2)

}

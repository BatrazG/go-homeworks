package main

import (
	"fmt"
	"sync"
)

func factorial(n int) int {
	if n == 0 {
		return 1
	}
	return n * factorial(n-1)
}

func main() {
	numbers := []int{2, 4, 6, 8, 10}

	// заранее выделяем память под результаты
	results := make([]int, len(numbers))

	var wg sync.WaitGroup
	wg.Add(len(numbers))

	for i, number := range numbers {
		// создаём локальные копии i и number для горутины
		i := i
		number := number

		go func() {
			defer wg.Done()

			result := factorial(number)
			results[i] = result

			// выводим сразу, как того просит условие
			fmt.Printf("Факториал %d = %d\n", number, result)
		}()
	}

	wg.Wait()

	fmt.Println("results:", results)
}

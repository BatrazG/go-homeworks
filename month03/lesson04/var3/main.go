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

// CalculateFactorial считает факториал и записывает результат по индексу i
func CalculateFactorial(n int, i int, results []int, wg *sync.WaitGroup) {
	defer wg.Done()

	result := factorial(n)
	results[i] = result

	fmt.Printf("Факториал %d = %d\n", n, result)
}

func main() {
	numbers := []int{2, 4, 6, 8, 10}

	results := make([]int, len(numbers))

	var wg sync.WaitGroup
	wg.Add(len(numbers))

	for i, number := range numbers {
		// важно: создаём локальные копии для замыкания
		i := i
		number := number

		go CalculateFactorial(number, i, results, &wg)
	}

	wg.Wait()

	fmt.Println("results:", results)
}

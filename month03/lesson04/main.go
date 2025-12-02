package main

import (
	"fmt"
	"sync"
)

func factorial(n int) int {
	if n == 0 || n == 1 {
		return 1
	}
	return n * factorial(n-1)
}

func main() {
	numbers := []int{5, 2, 10, 7}
	results := []int{}
	var mu sync.Mutex

	var wg sync.WaitGroup

	for _, number := range numbers {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			result := factorial(n)

			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(number)
	}

	wg.Wait()

	for i, r := range results {
		fmt.Printf("%d! = %d\n", numbers[i], r)
	}
}

package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Объявляем контракт: кто-то, кто умеет проверять сайты
type SiteChecker interface {
	Check(url string)
}

// Реализация 1: реальный HTTP чекер
type httpChecker struct{}

// Метод Check привязан к структуре httpChecker
func (h httpChecker) Check(url string) {
	// Код из старой функции checkUrl
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("ERR %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	duration := time.Since(start)
	fmt.Printf("[%d]%s - %v\n", resp.StatusCode, url, duration)
}

// Функция процессор, которая принимает интерфейс
func processURLs(checker SiteChecker, urls []string) {
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			checker.Check(u)
		}(url)
	}
	wg.Wait()
}

func main() {
	urls := []string{
		"https://google.com",
		"https://ya.ru",
		"https://github.com",
		"https://golang.org",
		"https://stackoverflow.com",
	}
	start := time.Now()

	// Создаем экземпляр нашей реализации
	myChecker := httpChecker{}

	// Передаем его в функцию, которая ожидает интерфейс
	processURLs(myChecker, urls)

	fmt.Printf("Total time: %v\n", time.Since(start))
}

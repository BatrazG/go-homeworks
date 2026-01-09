package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type Inspector interface {
	Inspect(url string)
}

type RealInspector struct{}

func (r RealInspector) Inspect(url string) {
	// Создаем клиент с таймаутом для защиты от зависаний
	// 10 секунд - баланс между терпением пользователя и отзывчивостью системы
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	// Используем кастомный клиент вместо http.Get()
	resp, err := client.Get(url)

	if err != nil {
		fmt.Printf("Ошибка запроса: %v\n", err)
		return
	}
	// Обязательно закрываем тело, даже при ошибках
	defer resp.Body.Close()
	// Проверяем статус-код ответа: успешными считаем только 2xx
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Printf("Неожиданный статус: %d %s\n",
			resp.StatusCode, http.StatusText(resp.StatusCode))
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Ошибка чтения тела запроса: %v\n", err)
		return
	}

	// Добавим ыременные метки в логирование
	fmt.Printf("[%s] URL: %s | Размер: %d байт\n",
		time.Now().Format("15:04:05"), url, len(body))
}

func main() {
	urls := []string{
		"https://mates-web.ru/",
		"https://goolgle.com",
		"https://journal.top-academy.ru",
		"https://ironau.ru",
		"https://github.com",
	}
	inspector := RealInspector{}
	var wg sync.WaitGroup

	// Паттерин семафор для ограничения параллелизма
	sem := make(chan struct{}, 5) // Максимум 5 одновременных запросов

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			// recover() ловит панику только в своей горутине,
			// не давая упасть всей программе из-за ошибки в одном URL
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Паника в горутине для %s: %v\n", u, r)
				}
			}()
			sem <- struct{}{}        // Занимаем слот
			defer func() { <-sem }() // Освобождаем слот
			defer wg.Done()
			inspector.Inspect(u)
		}(url)
	}

	wg.Wait()
}

package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func getPageSize(url string) (int, error) {
	// Создаем клиент с таймаутом для защиты от зависаний
	// 10 секунд - баланс между терпением пользователя и отзывчивостью системы
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	// Используем кастомный клиент вместо http.Get()
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("неожиданный статус: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	return len(body), nil
}

func main() {
	size, err := getPageSize("https://example.com")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println("Размер body:", size)

}

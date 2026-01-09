package main

import (
	"fmt"
	"io"
	"net/http"
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
	// Проверяем статус-код ответа: успешными считаем только 2xx
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Printf("Неожиданный статус: %d %s\n",
			resp.StatusCode, http.StatusText(resp.StatusCode))
		return
	}

	// Обязательно закрываем тело, даже при ошибках
	defer resp.Body.Close()

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
	inspector := RealInspector{}
	inspector.Inspect("https://mates-web.ru/")
}

package main

import (
	"fmt"
	"time"
)

// Функции-генераторы сообщений
func newsFeed(ch chan<- string) {
	for i := 1; i <= 9; i++ {
		freshNew := fmt.Sprintf("Новость: %d", i)
		ch <- freshNew
		time.Sleep(1 * time.Second)
	}
	close(ch)
}

func socialMedia(ch chan<- string) {
	for i := 1; i <= 3; i++ {
		media := fmt.Sprintf("Соцсети: %d", i)
		ch <- media
		time.Sleep(3 * time.Second)
	}
	close(ch)
}

func main() {
	news := make(chan string)
	social := make(chan string)
	// Глобальный таймаут на всю программу.
	// Важно: создаем ОДИН раз, до цикла, иначе на каждой итерации
	// таймер будет сбрасываться, и таймаут никогда не наступит
	timeout := time.After(5 * time.Second)

	go newsFeed(news)
	go socialMedia(social)

	for {
		if news == nil && social == nil {
			fmt.Println("\nКаналы закрыты, выключаемся")
			return
		}

		select {
		case msg1, ok := <-news:
			if !ok {
				news = nil
				continue
			}
			fmt.Print("\n")
			fmt.Println(msg1)
		case msg2, ok := <-social:
			if !ok {
				social = nil
				continue
			}
			fmt.Print("\n")
			fmt.Println(msg2)
		case <-timeout:
			fmt.Println("\nВремя вышло, выключаемся")
			return
		default:
			fmt.Print(".")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

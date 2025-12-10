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
	stopChan := make(chan struct{})

	go newsFeed(news)
	go socialMedia(social)
	go func() {
		time.Sleep(5 * time.Second)
		close(stopChan)
	}()

stop:
	for {
		if news == nil && social == nil {
			fmt.Println("\nКаналы закрыты, выключаемся")
			return
		}
		select {
		case msg1, ok := <-news:
			if !ok {
				// присваиваем nil, чтобы select больше не выбирал этот case
				news = nil
				continue
			}
			fmt.Print("\n")
			fmt.Println(msg1)
		case msg2, ok := <-social:
			if !ok {
				// присваиваем nil, чтобы select больше не выбирал этот case
				social = nil
				continue
			}
			fmt.Print("\n")
			fmt.Println(msg2)
		case <-stopChan:
			fmt.Println("\nВремя вышло, выключаемся")
			break stop
		default:
			fmt.Print(".")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

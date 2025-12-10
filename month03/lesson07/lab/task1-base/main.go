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

	go newsFeed(news)
	go socialMedia(social)

	for {
		if news == nil && social == nil {
			break
		}
		select {
		case msg1, ok := <-news:
			if !ok {
				news = nil
				continue
			}
			fmt.Println(msg1)
		case msg2, ok := <-social:
			if !ok {
				social = nil
				continue
			}
			fmt.Println(msg2)
		}

	}
}

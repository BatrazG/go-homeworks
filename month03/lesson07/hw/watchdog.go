package main

import (
	"fmt"
	"time"
)

func worker(ch chan<- string) {
	for i := 0; i < 5; i++ {
		defer close(ch)
		time.Sleep(1 * time.Second)
		ch <- "ping"
	}
	// "Смерть" воркера — больше ничего не шлём
	time.Sleep(3 * time.Second)
}

func main() {
	heartbeat := make(chan string)

	go worker(heartbeat)

	for {
		timeout := time.After(2 * time.Second)
		select {
		case msg := <-heartbeat:
			fmt.Println(msg)
		case <-timeout:
			fmt.Println("CRITICAL: Service dead")
			return
		}
	}
}

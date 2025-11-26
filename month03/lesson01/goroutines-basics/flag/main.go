package main

import (
	"flag"
	"fmt"
	"time"
)

const (
	perGoroutineWait    = 10 * time.Millisecond
	messagePerGoroutine = 3
)

func printMessage(id int) {
	for i := 0; i < messagePerGoroutine; i++ {
		fmt.Printf("[goroutine %d] message %d\n", id, i)
	}
}

func main() {
	n := flag.Int("n", 10, "сколько горутин запустить")

	flag.Parse()

	for i := 0; i < *n; i++ {
		go printMessage(i)
	}

	// ВРЕМЕННОЕ решение: ждём "примерно достаточно" времени, чтобы горутины успели отработать.
	// Это ненадёжно: если работа горутин займёт дольше, часть сообщений может не успеть напечататься.
	// В реальном коде нужно использовать явную синхронизацию (например, WaitGroup), а не Sleep.
	time.Sleep(time.Duration(*n) * perGoroutineWait)
}

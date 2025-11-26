/*
main принимает число N (количество горутин) из аргумента командной строки или флага;
запускает N горутин;
каждая горутина несколько раз печатает сообщения в формате:
[goroutine X] message Y
где `X` — номер горутины, `Y` — номер сообщения.
в конце main делает «грубое ожидание» через time.Sleep (с запасом, чтобы почти всегда «успевали» все горутины).
*/
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	perGorutineWait     = 10 * time.Millisecond
	messagePerGoroutine = 3
)

func printMessage(id int) {
	for i := 0; i < messagePerGoroutine; i++ {
		fmt.Printf("[goroutine %d] message %d\n", id, i)
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("ожидался ровно один аргумент (кол-во горутин), получено: %d\n", len(os.Args)-1)
	}
	numOfGoroutines, err := strconv.Atoi(os.Args[1])
	if err != nil || numOfGoroutines <= 0 {
		log.Fatalf("аргумент должен быть положительным числом, получено: %q\n", os.Args[1])
	}
	for i := 0; i < numOfGoroutines; i++ {
		go printMessage(i)
	}
	time.Sleep(time.Duration(numOfGoroutines) * perGorutineWait)
}

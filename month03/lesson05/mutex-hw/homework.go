package main

import (
	"fmt"
	"strings"
	"sync"
)

func ConcurrentWordCount(sentences []string) map[string]int {
	wordCount := make(map[string]int)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, sentence := range sentences {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()

			local := make(map[string]int)
			for _, w := range strings.Fields(s) {
				local[w]++
			}

			mu.Lock()
			for w, c := range local {
				wordCount[w] += c
			}
			mu.Unlock()
		}(sentence)
	}
	wg.Wait()
	return wordCount
}

func main() {
	text := []string{
		"quick brown fox",
		"lazy dog",
		"quick brown fox jumps",
		"jumps over lazy dog",
	}

	result := ConcurrentWordCount(text)
	fmt.Println(result)

}

/*
Мини-эссе: Mutex vs Каналы
Почему Mutex? В этой задаче нужна общая map, куда пишут все горутины.
 Mutex проще и эффективнее - одна блокировка на запись.

Когда каналы лучше? Каналы удобнее для передачи данных между горутинами без shared memory.
Здесь можно было бы создать канал, куда каждая горутина отправляет map со своими словами,
а одна горутина-агрегатор собирает результат - это избавило бы от блокировок, но добавило бы сложности.
*/

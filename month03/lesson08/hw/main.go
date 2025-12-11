package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// Структура для результатов
type Results struct {
	FileName string
	Lines    int
}

// Шаг 1-2. Источник и фильтр. Вариант с объединением
func filterSource(files []string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for _, f := range files {
			if strings.HasSuffix(f, ".txt") {
				out <- f
			}
		}
	}()
	return out
}

// Шаг 3. Обработка файла (воркер)
func fileWorker(id int, in <-chan string) <-chan Results {
	out := make(chan Results)
	go func() {
		defer close(out)
		for f := range in {
			fmt.Printf("Worker %d processing %s...\n", id, f)
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(300))) // Как будто считаем строки
			out <- Results{
				FileName: f,
				Lines:    rand.Intn(1000),
			}
		}
	}()
	return out
}

// Шаг 4. Слияние каналов(Fan-in)
func merge(cs ...<-chan Results) <-chan Results {
	out := make(chan Results)
	var wg sync.WaitGroup

	output := func(c <-chan Results) {
		defer wg.Done()
		for n := range c {
			out <- n
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	files := []string{"log1.txt", "image.jpg", "log2.txt", "data.csv", "readme.txt"}

	// 1. Pipeline
	filesChan := filterSource(files)

	// 2. Fan-out Запускаем 3 независимых воркера
	w1 := fileWorker(1, filesChan)
	w2 := fileWorker(2, filesChan)
	w3 := fileWorker(3, filesChan)

	// 3. Fan-in
	resultsChan := merge(w1, w2, w3)

	// 4. Получаем и выводим результаты
	totalLines := 0
	for res := range resultsChan {
		fmt.Printf("Done: %s, (%d lines)\n", res.FileName, res.Lines)
		totalLines += res.Lines
	}
	fmt.Printf("Total lines counted: %d", totalLines)
}

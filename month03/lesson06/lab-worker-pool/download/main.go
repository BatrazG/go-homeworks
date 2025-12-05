package main

import (
	"fmt"
	"sync"
	"time"
)

// Job представляет задачу скачивания файла
type Job struct {
	ID  int
	URL string
}

// Result содержит результат выполнения задачи
type Result struct {
	JobID    int
	Status   string // Статус скачивания
	WorkerID int
}

// worker обрабатывает задачи из канала jobs
func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Воркер %d начал смену\n", id)

	for job := range jobs {
		// Имитация скачивания файла
		time.Sleep(500 * time.Millisecond)

		status := fmt.Sprintf("Скачано с %s", job.URL)

		results <- Result{
			JobID:    job.ID,
			Status:   status,
			WorkerID: id,
		}
	}

	fmt.Printf("Воркер %d закончил смену\n", id)
}

func main() {
	// Каналы для задач и результатов
	jobs := make(chan Job)
	results := make(chan Result, 5) // Буфер для результатов

	var wg sync.WaitGroup

	// Параметры
	numWorkers := 3
	numJobs := 5 // По условию - 5 сайтов

	// Список тестовых URL
	urls := []string{
		"https://example.com/file1.zip",
		"https://test.org/archive.tar.gz",
		"https://demo.net/package.deb",
		"https://sample.io/installer.exe",
		"https://mock.dev/update.bin",
	}

	// 1. Запускаем воркеров
	fmt.Printf("Запускаем %d воркеров\n", numWorkers)
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}

	// 2. Отправляем задачи (в отдельной горутине)
	go func() {
		fmt.Println("Начинаем отправку задач")
		for j := 0; j < numJobs; j++ {
			jobs <- Job{
				ID:  j + 1,
				URL: urls[j],
			}
			time.Sleep(100 * time.Millisecond)
		}
		close(jobs) // Закрываем после отправки всех задач
		fmt.Println("Все задачи отправлены")
	}()

	// 3. Закрываем results после завершения всех воркеров
	go func() {
		wg.Wait()
		close(results)
	}()

	// 4. Читаем результаты
	fmt.Println("Ожидаем результаты...")
	fmt.Println()
	for res := range results {
		fmt.Printf("Задача %d: %s (Воркер %d)\n",
			res.JobID, res.Status, res.WorkerID)
	}

	fmt.Println("\nВсе задачи выполнены. Работа завершена.")
}

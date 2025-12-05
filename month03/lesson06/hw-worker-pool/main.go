package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Job - задача на обработку
type Job struct {
	ID     int
	Number int
}

// Result - результат вычислений
type Result struct {
	JobID     int
	InputNum  int
	Square    int // Площадь квадрата
	Perimeter int // Периметр квадрата
	WorkerID  int
}

// worker - обрабатывает ВСЕ задачи из канала jobs
func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	// ✅ Цикл для обработки ВСЕХ задач
	for job := range jobs {
		// Имитация сложных вычислений
		time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)

		square := job.Number * job.Number // Площадь
		perimeter := 4 * job.Number       // Периметр

		results <- Result{
			JobID:     job.ID,
			InputNum:  job.Number,
			Square:    square,
			Perimeter: perimeter,
			WorkerID:  id,
		}
	}

	fmt.Printf("🛑 Воркер %d завершил работу\n", id)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Входные данные
	inputs := []int{1, 5, 12, 5, 3, 8, 9}

	// Каналы (без буфера для простоты, но можно len(inputs))
	jobs := make(chan Job)
	results := make(chan Result)

	var wg sync.WaitGroup
	const numWorkers = 3 // ✅ По условию

	// 1️⃣ Запускаем 3 воркера
	fmt.Println("🚀 Запускаем воркеров...")
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}

	// 2️⃣ Отправляем задачи (отдельная горутина)
	go func() {
		fmt.Println("📤 Отправка задач началась")
		fmt.Println()
		for i, num := range inputs {
			jobs <- Job{ID: i + 1, Number: num}
		}
		close(jobs) // ✅ Закрываем после отправки всех задач
		fmt.Println("✅ Все задачи отправлены")
	}()

	// 3️⃣ Закрываем results после завершения воркеров
	go func() {
		wg.Wait()      // Ждём завершения всех воркеров
		close(results) // Закрываем канал результатов
	}()

	// 4️⃣ Читаем результаты (main не знает количество заранее)
	fmt.Println("📥 Получаем результаты:")
	fmt.Println()
	for res := range results {
		fmt.Printf("Задача #%d | Сторона: %d → Площадь: %d, Периметр: %d | Воркер: %d\n",
			res.JobID, res.InputNum, res.Square, res.Perimeter, res.WorkerID)
	}

	fmt.Println("\n🎉 Все задачи обработаны!")
}

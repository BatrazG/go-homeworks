package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// Config хранит конфигурацию программы
type Config struct {
	Timeout       time.Duration // Таймаут для HTTP запросов
	MaxConcurrent int           // Максимальное количество одновременных запросов
	OutputFile    string        // Файл для записи результатов
	HTTPMethod    string        // HTTP метод для запросов
}

// Result хранит результат проверки одного URL
type Result struct {
	URL      string        // Проверяемый URL
	Status   int           // HTTP статус-код (0 если ошибка)
	Duration time.Duration // Время выполнения запроса
	Error    error         // Ошибка (если есть)
}

// checkURL выполняет HTTP запрос к указанному URL
func checkURL(url string, config Config, sem chan struct{}, wg *sync.WaitGroup, results chan<- Result) {
	defer wg.Done() // Уменьшаем счетчик WaitGroup при завершении

	// Захватываем слот в семафоре для ограничения одновременных запросов
	sem <- struct{}{}
	defer func() { <-sem }() // Освобождаем слот при выходе из функции

	start := time.Now() // Засекаем время начала запроса

	// Создаем HTTP клиент с указанным таймаутом
	client := &http.Client{
		Timeout: config.Timeout,
	}

	// Создаем HTTP запрос с указанным методом
	req, err := http.NewRequest(config.HTTPMethod, url, nil)
	if err != nil {
		results <- Result{
			URL:      url,
			Status:   0,
			Duration: time.Since(start),
			Error:    fmt.Errorf("ошибка создания запроса: %w", err),
		}
		return
	}

	// Устанавливаем User-Agent для идентификации нашего запроса
	req.Header.Set("User-Agent", "URL-Checker/1.0")

	// Выполняем HTTP запрос
	resp, err := client.Do(req)
	if err != nil {
		results <- Result{
			URL:      url,
			Status:   0,
			Duration: time.Since(start),
			Error:    fmt.Errorf("ошибка запроса: %w", err),
		}
		return
	}
	defer resp.Body.Close() // Важно: закрываем тело ответа для освобождения ресурсов

	// Читаем тело ответа до конца (освобождает соединение для повторного использования)
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		results <- Result{
			URL:      url,
			Status:   resp.StatusCode,
			Duration: time.Since(start),
			Error:    fmt.Errorf("ошибка чтения тела ответа: %w", err),
		}
		return
	}

	// Возвращаем успешный результат
	results <- Result{
		URL:      url,
		Status:   resp.StatusCode,
		Duration: time.Since(start),
		Error:    nil,
	}
}

// readURLsFromFile читает URL из текстового файла
func readURLsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла %s: %w", filename, err)
	}
	defer file.Close() // Закрываем файл при выходе из функции

	var urls []string
	scanner := bufio.NewScanner(file)

	// Читаем файл построчно
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Пропускаем пустые строки и комментарии
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Добавляем URL в список
		urls = append(urls, line)
	}

	// Проверяем ошибки сканера
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения файла %s: %w", filename, err)
	}

	return urls, nil
}

// checkURLsConcurrent запускает параллельную проверку URL
func checkURLsConcurrent(urls []string, config Config) []Result {
	// Создаем семафор для ограничения количества одновременных запросов
	sem := make(chan struct{}, config.MaxConcurrent)

	// Канал для сбора результатов
	resultsChan := make(chan Result, len(urls))

	// WaitGroup для ожидания завершения всех горутин
	var wg sync.WaitGroup

	// Запускаем горутины для каждого URL
	for _, url := range urls {
		wg.Add(1)
		go checkURL(url, config, sem, &wg, resultsChan)
	}

	// Запускаем горутину, которая закроет канал результатов после завершения всех запросов
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Собираем результаты из канала
	var results []Result
	for result := range resultsChan {
		results = append(results, result)
	}

	return results
}

// printResults выводит результаты проверки
func printResults(results []Result, config Config) error {
	// Создаем multiWriter для вывода одновременно в консоль и файл
	var writers []io.Writer
	writers = append(writers, os.Stdout)

	// Если указан файл для вывода, открываем его
	if config.OutputFile != "" {
		file, err := os.Create(config.OutputFile)
		if err != nil {
			return fmt.Errorf("ошибка создания файла отчета: %w", err)
		}
		defer file.Close()
		writers = append(writers, file)
	}

	// Создаем multiWriter
	multiWriter := io.MultiWriter(writers...)

	// Выводим заголовок отчета
	fmt.Fprintf(multiWriter, "\n=== ОТЧЕТ ПО ПРОВЕРКЕ URL ===\n")
	fmt.Fprintf(multiWriter, "Метод HTTP: %s\n", config.HTTPMethod)
	fmt.Fprintf(multiWriter, "Таймаут: %v\n", config.Timeout)
	fmt.Fprintf(multiWriter, "Макс. одновременных запросов: %d\n\n", config.MaxConcurrent)

	// Счетчики для статистики
	successCount := 0
	errorCount := 0

	// Выводим результаты для каждого URL
	for _, result := range results {
		if result.Error != nil {
			// Выводим информацию об ошибке
			fmt.Fprintf(multiWriter, "[ERROR] %s (%v) - %v\n",
				result.URL, result.Duration.Round(time.Millisecond), result.Error)
			errorCount++
		} else {
			// Выводим успешный результат
			fmt.Fprintf(multiWriter, "[%d] %s (%v)\n",
				result.Status, result.URL, result.Duration.Round(time.Millisecond))
			successCount++
		}
	}

	// Выводим сводку
	fmt.Fprintf(multiWriter, "\n=== СВОДКА ===\n")
	fmt.Fprintf(multiWriter, "Всего URL: %d\n", len(results))
	fmt.Fprintf(multiWriter, "Успешно: %d\n", successCount)
	fmt.Fprintf(multiWriter, "Ошибок: %d\n", errorCount)

	// Если есть файл вывода, сообщаем об этом
	if config.OutputFile != "" {
		fmt.Fprintf(os.Stdout, "\nОтчет также сохранен в файл: %s\n", config.OutputFile)
	}

	return nil
}

func main() {
	// Определяем флаги командной строки
	var (
		timeoutStr    string
		maxConcurrent int
		outputFile    string
		httpMethod    string
	)

	// Устанавливаем значения по умолчанию и описания флагов
	flag.StringVar(&timeoutStr, "t", "10s", "Таймаут для HTTP запросов (например: 5s, 30s, 1m)")
	flag.IntVar(&maxConcurrent, "c", 10, "Максимальное количество одновременных запросов")
	flag.StringVar(&outputFile, "o", "", "Файл для записи результатов (опционально)")
	flag.StringVar(&httpMethod, "m", "GET", "HTTP метод для запросов (GET, HEAD, OPTIONS)")
	flag.Parse()

	// Парсим таймаут из строки
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		fmt.Printf("Ошибка парсинга таймаута '%s': %v\n", timeoutStr, err)
		fmt.Println("Используйте корректный формат: 5s, 30s, 1m, 2m30s и т.д.")
		os.Exit(1)
	}

	// Проверяем корректность HTTP метода
	validMethods := map[string]bool{
		"GET":     true,
		"HEAD":    true,
		"OPTIONS": true,
		"POST":    false, // POST требует тела запроса
		"PUT":     false, // PUT требует тела запроса
		"DELETE":  true,
	}

	if !validMethods[httpMethod] {
		fmt.Printf("Предупреждение: метод '%s' может требовать тела запроса\n", httpMethod)
		fmt.Println("Рекомендуемые методы для проверки доступности: GET, HEAD, OPTIONS")
	}

	// Создаем конфигурацию
	config := Config{
		Timeout:       timeout,
		MaxConcurrent: maxConcurrent,
		OutputFile:    outputFile,
		HTTPMethod:    strings.ToUpper(httpMethod),
	}

	fmt.Printf("Запуск загрузчика с параметрами:\n")
	fmt.Printf("  Таймаут: %v\n", config.Timeout)
	fmt.Printf("  Макс. одновременных запросов: %d\n", config.MaxConcurrent)
	fmt.Printf("  HTTP метод: %s\n", config.HTTPMethod)
	if config.OutputFile != "" {
		fmt.Printf("  Файл отчета: %s\n", config.OutputFile)
	}
	fmt.Println()

	// Читаем URL из файла
	fmt.Println("Чтение URL из файла urls.txt...")
	urls, err := readURLsFromFile("urls.txt")
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		os.Exit(1)
	}

	// Проверяем, что файл не пуст
	if len(urls) == 0 {
		fmt.Println("Предупреждение: файл urls.txt пуст или не содержит валидных URL")
		fmt.Println("Добавьте URL в формате: одна строка - один URL")
		fmt.Println("Комментарии начинаются с #")
		os.Exit(0)
	}

	fmt.Printf("Найдено %d URL для проверки\n", len(urls))
	fmt.Println("Начинаю параллельную проверку...")

	// Проверяем URL параллельно
	results := checkURLsConcurrent(urls, config)

	// Выводим результаты
	err = printResults(results, config)
	if err != nil {
		fmt.Printf("Ошибка при выводе результатов: %v\n", err)
		os.Exit(1)
	}
}

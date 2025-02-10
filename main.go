package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"task-meneger/pkg/storage"
	"task-meneger/pkg/storage/postgres"

	"github.com/joho/godotenv"
)

func main() {

	// Загружаем .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Не удалось загрузить .env файл, используем переменные окружения")
	}

	log.Println("DB_HOST:", os.Getenv("DB_HOST"))
	log.Println("DB_PORT:", os.Getenv("DB_PORT"))
	log.Println("DB_USER:", os.Getenv("DB_USER"))
	log.Println("DB_NAME:", os.Getenv("DB_NAME"))

	// Подключение к БД
	storage, err := postgres.New()
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer storage.Close() // Закрываем пул соединений при выходе

	// Приветствие
	fmt.Println("Добро пожаловать в Task Manager!")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\nВыберите действие:")
		fmt.Println("1. Посмотреть список задач")
		fmt.Println("2. Создать новую задачу")
		fmt.Println("3. Выйти")

		fmt.Print("Введите номер действия: ")
		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			printTasks(storage)
		case "2":
			createTask(scanner, storage)
		case "3":
			fmt.Println("Выход...")
			return
		default:
			fmt.Println("Некорректный ввод, попробуйте снова.")
		}
	}
}

// Функция для вывода списка задач
func printTasks(storage storage.Interface) {
	tasks, err := storage.Tasks(0, 0) // Получаем все задачи
	if err != nil {
		fmt.Println("Ошибка при получении списка задач:", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("Задач пока нет.")
		return
	}

	fmt.Println("\nСписок задач:")
	for _, task := range tasks {
		fmt.Printf("ID: %d | Title: %s | Content: %s\n", task.ID, task.Title, task.Content)
	}
}

// Функция для создания новой задачи
func createTask(scanner *bufio.Scanner, storage storage.Interface) {
	fmt.Print("Введите заголовок задачи: ")
	scanner.Scan()
	title := strings.TrimSpace(scanner.Text())

	fmt.Print("Введите описание задачи: ")
	scanner.Scan()
	content := strings.TrimSpace(scanner.Text())

	task := postgres.Task{
		Title:   title,
		Content: content,
	}

	id, err := storage.NewTask(task)
	if err != nil {
		fmt.Println("Ошибка при создании задачи:", err)
		return
	}

	fmt.Printf("Задача успешно создана! ID: %d\n", id)
}

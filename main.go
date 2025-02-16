package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"task-meneger/pkg/storage"
	"task-meneger/pkg/storage/postgres"

	"github.com/joho/godotenv"
)

func main() {

	// Загрузка переменных окружения из env
	err := godotenv.Load()
	if err != nil {
		log.Println("Не удалось загрузить .env файл, используем переменные окружения")
	}

	// Подключение к БД
	var storage storage.Interface
	storage, err = postgres.New()
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer storage.Close()

	// Приветствие и вывод меню в терминале
	fmt.Println("-------------------------------")
	fmt.Println("Добро пожаловать в Task Manager!")
	fmt.Println("-------------------------------")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n=============TASK==============")
		fmt.Println("1. Посмотреть список задач")
		fmt.Println("2. Создать новую задачу")
		fmt.Println("3. Обновить задачу")
		fmt.Println("4. Удалить задачу")
		fmt.Println("\n============LABELS=============")
		fmt.Println("5. Посмотреть список меток")
		fmt.Println("6. Создать новую метку")
		fmt.Println("\n============USERS==============")
		fmt.Println("7. Посмотреть список пользователей")
		fmt.Println("8. Создать нового пользователя")
		fmt.Println("\n============SEARCH=============")
		fmt.Println("9. Поиск задач по автору")

		fmt.Println("10. Выйти")

		fmt.Print("\nВведите номер действия: ")
		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			printTasks(storage)
			waitForEnter(scanner)
		case "2":
			createTask(scanner, storage)
			waitForEnter(scanner)
		case "3":
			updateTask(scanner, storage)
			waitForEnter(scanner)
		case "4":
			deleteTask(scanner, storage)
			waitForEnter(scanner)
		case "5":
			printLabels(storage)
			waitForEnter(scanner)
		case "6":
			createLabel(scanner, storage)
			waitForEnter(scanner)
		case "7":
			printUsers(storage)
			waitForEnter(scanner)
		case "8":
			createUser(scanner, storage)
			waitForEnter(scanner)
		case "9":
			getTasksByIdUser(scanner, storage)
			waitForEnter(scanner)

		case "10":
			fmt.Println("Выход...")
			return
		default:
			fmt.Println("\n🔴 Некорректный ввод, попробуйте снова.")
		}
	}
}

// Функция для возвращения в меню из выполненой функции
// Если убрать - то при выполнении любой функции будет появляться меню
func waitForEnter(scanner *bufio.Scanner) {
	fmt.Println("\n🔙 Нажмите Enter, чтобы вернуться в главное меню")
	scanner.Scan()
}

// Функция для вывода списка задач
func printTasks(storage storage.Interface) {
	tasks, err := storage.Tasks(0, 0) // Получаем все задачи
	if err != nil {
		fmt.Println("-------------------------------")
		fmt.Println("\n🔴 Ошибка при получении списка задач:", err)
		fmt.Println("-------------------------------")
		return
	}

	if len(tasks) == 0 {
		fmt.Println("-------------------------------")
		fmt.Println("\n⚠️  Задач пока нет.")
		fmt.Println("-------------------------------")
		return
	}

	fmt.Println("\n📋 Список задач:")
	for _, task := range tasks {
		fmt.Println("-------------------------------")
		fmt.Printf("🆔 ID: %d\n📌 Заголовок: %s\n📝 Описание: %s\n👤 Автор: %d\n🎯 Исполнитель: %d\n",
			task.ID, task.Title, task.Content, task.AuthorID, task.AssignedID)
	}
	fmt.Println("-------------------------------")
}

// Функция для создания новой задачи
func createTask(scanner *bufio.Scanner, storage storage.Interface) {
	fmt.Println("-------------------------------")
	fmt.Print("\n📌 Введите заголовок задачи: ")
	fmt.Println("-------------------------------")
	scanner.Scan()
	title := strings.TrimSpace(scanner.Text())

	fmt.Print("\n📝 Введите описание задачи: ")
	fmt.Println("-------------------------------")
	scanner.Scan()
	content := strings.TrimSpace(scanner.Text())

	fmt.Print("\n👤 Введите ID автора: ")
	fmt.Println("-------------------------------")
	scanner.Scan()
	authorID, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		fmt.Println("\n🔴 Ошибка: Некорректный ID автора")
		fmt.Println("-------------------------------")
		return
	}

	fmt.Println("-------------------------------")
	fmt.Print("\n🎯 Введите ID исполнителя: ")
	fmt.Println("-------------------------------")
	scanner.Scan()
	assignedID, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		fmt.Println("-------------------------------")
		fmt.Println("🔴 Ошибка: Некорректный ID исполнителя")
		fmt.Println("-------------------------------")
		return
	}

	// Ввод меток (можно несколько через запятую)
	fmt.Println("-------------------------------")
	fmt.Print("\n🏷️  Введите ID меток через запятую (или оставьте пустым): ")
	fmt.Println("-------------------------------")
	scanner.Scan()
	labelInput := strings.TrimSpace(scanner.Text())

	var labelIDs []int
	if labelInput != "" {
		labelStrs := strings.Split(labelInput, ",")
		for _, str := range labelStrs {
			labelID, err := strconv.Atoi(strings.TrimSpace(str))
			if err != nil {
				fmt.Println("-------------------------------")
				fmt.Println("🔴 Ошибка: Некорректный ID метки")
				fmt.Println("-------------------------------")
				return
			}
			labelIDs = append(labelIDs, labelID)
		}
	}

	// Создаём задачу
	task := postgres.Task{
		Title:      title,
		Content:    content,
		AuthorID:   authorID,
		AssignedID: assignedID,
	}

	id, err := storage.NewTask(task, labelIDs)
	if err != nil {
		fmt.Println("-------------------------------")
		fmt.Println("\n🔴 Ошибка при создании задачи:", err)
		fmt.Println("-------------------------------")
		return
	}

	fmt.Println("-------------------------------")
	fmt.Printf("\n✅ Задача успешно создана! ID: %d\n", id)
	fmt.Println("-------------------------------")
}

// Функция для обновления задачи
func updateTask(scanner *bufio.Scanner, storage storage.Interface) {
	fmt.Println("-------------------------------")
	fmt.Print("\n🆔 Введите ID задачи для обновления: ")
	fmt.Println("-------------------------------")
	scanner.Scan()
	taskID, _ := strconv.Atoi(scanner.Text()) // Преобразуем ввод в число

	fmt.Print("\n📌 Введите новый заголовок задачи: ")
	fmt.Println("-------------------------------")
	scanner.Scan()
	title := strings.TrimSpace(scanner.Text())

	fmt.Print("\n📝 Введите новое описание задачи: ")
	fmt.Println("-------------------------------")
	scanner.Scan()
	content := strings.TrimSpace(scanner.Text())

	fmt.Print("\n👤 Введите ID автора: ")
	fmt.Println("-------------------------------")
	scanner.Scan()
	authorID, _ := strconv.Atoi(scanner.Text())

	fmt.Print("\n🎯 Введите ID исполнителя: ")
	fmt.Println("-------------------------------")
	scanner.Scan()
	assignedID, _ := strconv.Atoi(scanner.Text())

	task := postgres.Task{
		ID:         taskID,
		Title:      title,
		Content:    content,
		AuthorID:   authorID,
		AssignedID: assignedID,
	}

	err := storage.UpdateTask(task)
	if err != nil {
		fmt.Println("-------------------------------")
		fmt.Println("\n🔴 Ошибка при обновлении задачи:", err)
		fmt.Println("-------------------------------")
		return
	}

	fmt.Println("\n✅ Задача успешно обновлена!")
	fmt.Println("-------------------------------")
}

// Функция для удаления задачи
func deleteTask(scanner *bufio.Scanner, storage storage.Interface) {
	fmt.Print("\n🆔 Введите ID задачи для удаления: ")
	scanner.Scan()
	taskID, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		fmt.Println("\n❌ Ошибка: Некорректный ID")
		return
	}

	err = storage.DeleteTask(taskID)
	if err != nil {
		fmt.Println("\n🔴 Ошибка при удалении задачи:", err)
		return
	}

	fmt.Println("\n✅ Задача успешно удалена!")
	fmt.Println("-------------------------------")
}

// Функция для вывода всех меток
func printLabels(storage storage.Interface) {
	labels, err := storage.Labels()
	if err != nil {
		fmt.Println("-------------------------------")
		fmt.Println("\n🔴 Ошибка при получении списка меток:", err)
		fmt.Println("-------------------------------")
		return
	}

	if len(labels) == 0 {
		fmt.Println("-------------------------------")
		fmt.Println("\n⚠️  Метки отсутствуют.")
		fmt.Println("-------------------------------")
		return
	}

	fmt.Println("-------------------------------")
	fmt.Println("\n🏷️  Список меток:")
	for _, label := range labels {
		fmt.Printf("ID: %d | Название: %s\n", label.ID, label.Name)
	}
	fmt.Println("-------------------------------")
}

// Функция для создания новой метки
func createLabel(scanner *bufio.Scanner, storage storage.Interface) {
	fmt.Println("-------------------------------")
	fmt.Print("\n🏷️  Введите название метки: ")
	fmt.Println("-------------------------------")
	scanner.Scan()

	name := strings.TrimSpace(scanner.Text())

	label := postgres.Label{
		Name: name,
	}

	id, err := storage.NewLabel(label)
	if err != nil {
		fmt.Println("\n 🔴 Ошибка при создании метки: ", err)
		fmt.Println("-------------------------------")
		return
	}

	fmt.Printf("\n✅ Метка успешно создана! ID: %d\n", id)
	fmt.Println("-------------------------------")

}

// Функция для вывода пользователей
func printUsers(storage storage.Interface) {
	users, err := storage.Users()
	if err != nil {
		fmt.Println("-------------------------------")
		fmt.Println("\n🔴 Ошибка при получении списка пользователей:", err)
		fmt.Println("-------------------------------")
		return
	}

	if len(users) == 0 {
		fmt.Println("-------------------------------")
		fmt.Println("\n⚠️  Пользователи отсутствуют.")
		fmt.Println("-------------------------------")
		return
	}

	fmt.Println("-------------------------------")
	fmt.Println("\n👤 Список пользователей:")
	for _, user := range users {
		fmt.Printf("ID: %d | Имя: %s\n", user.ID, user.Name)
	}
	fmt.Println("-------------------------------")
}

// Функция для созданя нового пользователя
func createUser(scanner *bufio.Scanner, storage storage.Interface) {
	fmt.Println("-------------------------------")
	fmt.Print("\n👤 Введите имя пользователя: ")
	fmt.Println("-------------------------------")
	scanner.Scan()

	name := strings.TrimSpace(scanner.Text())

	user := postgres.User{
		Name: name,
	}

	id, err := storage.NewUser(user)
	if err != nil {
		fmt.Println("\n🔴 Ошибка при создании пользователя: ", err)
		fmt.Println("-------------------------------")
		return
	}

	fmt.Printf("\n👤 Пользователь успешно создан! ID: %d\n", id)
	fmt.Println("-------------------------------")
}

// Функция для поиска задачи по автору
func getTasksByIdUser(scanner *bufio.Scanner, storage storage.Interface) {
	fmt.Println("-------------------------------")
	fmt.Println("\n👤 Введите ID автора: ")
	fmt.Println("-------------------------------")
	scanner.Scan()
	authorID, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("\n🔴 Ошибка: некорректный ID")
		fmt.Println("-------------------------------")
		return
	}

	tasks, err := storage.GetTasksByAuthor(authorID)
	if err != nil {
		fmt.Println("\n🔴 Ошибка при получении списка задач:", err)
		fmt.Println("-------------------------------")
		return
	}

	if len(tasks) == 0 {
		fmt.Println("\n⚠️  У данного автора пока нет задач.")
		fmt.Println("-------------------------------")
		return
	}

	fmt.Println("\nСписок задач автора:")
	for _, task := range tasks {
		fmt.Printf("🆔 ID: %d | 📌 Заголовок: %s | 📝 Описание: %s\n",
			task.ID, task.Title, task.Content)
	}
	fmt.Println("-------------------------------")
}

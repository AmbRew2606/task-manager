package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Хранилище данных.
type Storage struct {
	db *pgxpool.Pool //ПУЛ СОЕДИНЕНИЙ
}

// Функция New - подключение к БД
func New() (*Storage, error) {

	// err := godotenv.Load("../../../.env")
	// if err != nil {
	// 	log.Println("Не удалось загрузить .env файл, используем переменные окружения")
	// }

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// log.Println("DB_HOST:", os.Getenv("DB_HOST"))
	// log.Println("DB_PORT:", os.Getenv("DB_PORT"))
	// log.Println("DB_USER:", os.Getenv("DB_USER"))
	// log.Println("DB_NAME:", os.Getenv("DB_NAME"))

	// host := "localhost"
	// port := "5433"
	// user := "postgres"
	// password := "2606"
	// dbname := "db_task_manager"

	// для теста енв файла, если проблемы
	// log.Printf("DB_HOST=%s DB_PORT=%s DB_USER=%s DB_NAME=%s", host, port, user, dbname)

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname)

	dbpool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	return &Storage{db: dbpool}, nil
}

// Закрытие соединения с БД
func (s *Storage) Close() {
	s.db.Close()
}

// Задача.
type Task struct {
	ID         int
	Opened     int64
	Closed     int64
	AuthorID   int
	AssignedID int
	Title      string
	Content    string
}

// Метка.
type Label struct {
	ID   int
	Name string
}

// Пользователь.
type User struct {
	ID   int
	Name string
}

// Tasks возвращает список задач из БД.
func (s *Storage) Tasks(taskID, authorID int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
		WHERE
			($1 = 0 OR id = $1) AND
			($2 = 0 OR author_id = $2)
		ORDER BY id;
	`,
		taskID,
		authorID,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return tasks, rows.Err()
}

// NewTask создаёт новую задачу и возвращает её id.
func (s *Storage) NewTask(t Task, labelIDs []int) (int, error) {
	var taskID int
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO tasks (title, content, author_id, assigned_id)
		VALUES ($1, $2, $3, $4) RETURNING id;
		`,
		t.Title,
		t.Content,
		t.AuthorID,
		t.AssignedID,
	).Scan(&taskID)
	// return taskID , err
	if err != nil {
		return 0, fmt.Errorf("ошибка при создании задачи: %w", err)
	}

	// 2. Добавляем связи с метками в tasks_labels
	for _, labelID := range labelIDs {
		_, err := s.db.Exec(context.Background(), `
			INSERT INTO tasks_labels (task_id, label_id)
			VALUES ($1, $2);
		`, taskID, labelID)
		if err != nil {
			return 0, fmt.Errorf("ошибка при добавлении метки: %w", err)
		}
	}

	return taskID, nil

}

// UpdateTask обновляет задачу по id.
func (s *Storage) UpdateTask(t Task) error {
	_, err := s.db.Exec(context.Background(), `
		UPDATE tasks 
		SET title = $1, content = $2, author_id = $3, assigned_id = $4
		WHERE id = $5;
		`,
		t.Title,
		t.Content,
		t.AuthorID,
		t.AssignedID,
		t.ID)

	if err != nil {
		return fmt.Errorf("ошибка при обновлении задачи: %w", err)
	}
	return nil
}

// DeleteTask удаляет задачу по id.
func (s *Storage) DeleteTask(taskID int) error {
	_, err := s.db.Exec(context.Background(), `
		DELETE FROM tasks WHERE id = $1;
	`, taskID)

	if err != nil {
		return fmt.Errorf("ошибка при удалении задачи: %w", err)
	}

	fmt.Printf("Задача с ID %d успешно удалена!\n", taskID)
	return nil
}

// Labels возвращает список меток из БД.
func (s *Storage) Labels() ([]Label, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT id, name FROM labels ORDER BY id;
	`)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении меток: %w", err)
	}
	defer rows.Close()

	var labels []Label

	for rows.Next() {
		var l Label
		if err := rows.Scan(&l.ID, &l.Name); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании метки: %w", err)
		}
		labels = append(labels, l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке строк: %w", err)
	}

	return labels, nil
}

// NewLabel создает новую метку и возвращает её id.
func (s *Storage) NewLabel(l Label) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO labels (name)
		VALUES ($1)
		RETURNING id;
	`, l.Name).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("ошибка при создании метки: %w", err)
	}
	return id, nil
}

// Users возвращает список пользователей из БД.
func (s *Storage) Users() ([]User, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT id, name FROM users ORDER BY id;
	`)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении меток: %w", err)
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании пользователей: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке строк: %w", err)
	}

	return users, nil
}

// Users создает нового пользователя и возвращает его id.
func (s *Storage) NewUser(u User) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO users (name)
		VALUES ($1)
		RETURNING id;
	`, u.Name).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("ошибка при создании пользователя: %w", err)
	}
	return id, nil
}

// GetTasksByAuthor возвращает список задач по id автора.
func (s *Storage) GetTasksByAuthor(authorID int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT id, title, content, author_id, assigned_id 
		FROM tasks WHERE author_id = $1;
	`, authorID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении задач: %w", err)
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var t Task
		err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Content,
			&t.AuthorID,
			&t.AssignedID,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, rows.Err()
}

package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Хранилище данных.
type Storage struct {
	db *pgxpool.Pool
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
func (s *Storage) NewTask(t Task) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO tasks (title, content)
		VALUES ($1, $2) RETURNING id;
		`,
		t.Title,
		t.Content,
	).Scan(&id)
	return id, err
}

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

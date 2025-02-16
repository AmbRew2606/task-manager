package memdb

import "task-meneger/pkg/storage/postgres"

type DB struct {
	tasks  []postgres.Task
	labels []postgres.Label
	users  []postgres.User
	nextID int
}

func New() *DB {
	return &DB{nextID: 1}
}

// Tasks — Получение списка задач
func (db *DB) Tasks(int, int) ([]postgres.Task, error) {
	return db.tasks, nil
}

// NewTask — Создание новой задачи
func (db *DB) NewTask(task postgres.Task, labels []int) (int, error) {
	task.ID = db.nextID
	db.nextID++
	db.tasks = append(db.tasks, task)
	return task.ID, nil
}

// UpdateTask — Обновление задачи
func (db *DB) UpdateTask(updatedTask postgres.Task) error {
	for i, t := range db.tasks {
		if t.ID == updatedTask.ID {
			db.tasks[i] = updatedTask
			return nil
		}
	}
	return nil // Можно вернуть ошибку, если задача не найдена
}

// DeleteTask — Удаление задачи
func (db *DB) DeleteTask(id int) error {
	for i, t := range db.tasks {
		if t.ID == id {
			db.tasks = append(db.tasks[:i], db.tasks[i+1:]...)
			return nil
		}
	}
	return nil // Можно вернуть ошибку, если задача не найдена
}

// Labels — Получение всех меток
func (db *DB) Labels() ([]postgres.Label, error) {
	return db.labels, nil
}

// NewLabel — Добавление новой метки
func (db *DB) NewLabel(label postgres.Label) (int, error) {
	label.ID = len(db.labels) + 1
	db.labels = append(db.labels, label)
	return label.ID, nil
}

// Users — Получение всех пользователей
func (db *DB) Users() ([]postgres.User, error) {
	return db.users, nil
}

// NewUser — Создание нового пользователя
func (db *DB) NewUser(user postgres.User) (int, error) {
	user.ID = len(db.users) + 1
	db.users = append(db.users, user)
	return user.ID, nil
}

// GetTasksByAuthor — Получение задач по автору
func (db *DB) GetTasksByAuthor(authorID int) ([]postgres.Task, error) {
	var result []postgres.Task
	for _, t := range db.tasks {
		if t.AuthorID == authorID {
			result = append(result, t)
		}
	}
	return result, nil
}

// Close — Закрытие "БД"
func (db *DB) Close() {
	// В памяти ничего закрывать не нужно, просто заглушка
}

package storage

import "task-meneger/pkg/storage/postgres"

//Интерфес БД
type Interface interface {
	Tasks(int, int) ([]postgres.Task, error)
	NewTask(postgres.Task) (int, error)
	Labels() ([]postgres.Label, error) // Добавляем Labels()
	NewLabel(postgres.Label) (int, error)
	Users() ([]postgres.User, error)
	NewUser(postgres.User) (int, error)
	Close() // Добавляем метод закрытия соединения
}

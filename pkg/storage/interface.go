package storage

import "task-meneger/pkg/storage/postgres"

//Интерфес БД
type Interface interface {
	//Tasks
	Tasks(int, int) ([]postgres.Task, error)
	NewTask(postgres.Task, []int) (int, error)
	UpdateTask(postgres.Task) error
	DeleteTask(int) error
	//Labels
	Labels() ([]postgres.Label, error)
	NewLabel(postgres.Label) (int, error)
	//Users
	Users() ([]postgres.User, error)
	NewUser(postgres.User) (int, error)
	//Search
	GetTasksByAuthor(int) ([]postgres.Task, error)
	Close() // для закрытия соединения с БД
}

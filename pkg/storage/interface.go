package storage

import "task-meneger/pkg/storage/postgres"

//Интерфес БД
type Interface interface {
	Tasks(int, int) ([]postgres.Task, error)
	NewTask(postgres.Task) (int, error)
}

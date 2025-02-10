package memdb

type DB []postgres.Task

func (db DB) Tasks(int, int) ([]postgres.Task, error) {
	return db, nil
}
func (db DB) NewTask(postgres.Task) (int, error) {
	return 0, nil
}

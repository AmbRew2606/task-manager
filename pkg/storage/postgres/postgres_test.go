package postgres

import (
	"reflect"
	"testing"
)

func TestDatabaseConnection(t *testing.T) {
	_, err := New()
	if err != nil {
		t.Fatalf("Ошибка подключения к БД: %v", err)
	}
}

func TestStorage_Tasks(t *testing.T) {
	type args struct {
		taskID   int
		authorID int
	}
	tests := []struct {
		name    string
		s       *Storage
		args    args
		want    []Task
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Tasks(tt.args.taskID, tt.args.authorID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Tasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.Tasks() = %v, want %v", got, tt.want)
			}
		})
	}
}

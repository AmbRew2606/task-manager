package storage

import (
	"testing"
)

func TestDatabaseConnection(t *testing.T) {
	_, err := New()
	if err != nil {
		t.Fatalf("Ошибка подключения к БД: %v", err)
	}
}

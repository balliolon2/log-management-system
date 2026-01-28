package job

import (
	"context"
	"testing"

	"log-management-backend/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestDeleteOldLogs(t *testing.T) {
	// Mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewLogRepository(db)

	// Expect Delete Query
	// DELETE FROM logs WHERE timestamp < NOW() - INTERVAL '1 day' * $1
	mock.ExpectExec("DELETE FROM logs WHERE timestamp < NOW() - INTERVAL '1 day' \\* \\$1").
		WithArgs(7).
		WillReturnResult(sqlmock.NewResult(0, 10)) // Deleted 10 rows

	count, err := repo.DeleteOldLogs(context.Background(), 7)
	if err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if count != 10 {
		t.Errorf("expected 10 affected rows, got %d", count)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

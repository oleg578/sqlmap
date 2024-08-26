package sql2json

import (
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func BenchmarkRowsToJson(b *testing.B) {
	// Initialize a mock database and rows
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"Column 1", "Column 2"})
	rows.AddRow("Dummy_1", 1)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	rs, err := db.Query("SELECT 1")
	if err != nil {
		b.Fatalf("an error '%s' occurred when querying the database", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = RowsToJson(rs)
		if err != nil {
			b.Fatalf("an error '%s' occurred during RowsToJson execution", err)
		}
	}
	b.StopTimer()
}

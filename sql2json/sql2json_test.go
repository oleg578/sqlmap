package sql2json

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brianvoe/gofakeit/v7"
	"testing"
)

func BenchmarkRowsToJson(b *testing.B) {
	// Initialize a mock database and rows
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"Id", "Product", "Price", "Qty", "NullData", "Date"})
	for i := 0; i < 1; i++ {
		rows.AddRow(
			i,
			gofakeit.Product().Name,
			gofakeit.Product().Price,
			gofakeit.IntRange(10, 1000),
			nil,
			gofakeit.Date().String())
	}

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

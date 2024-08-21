package sql2json_test

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"sql2json"
	"testing"
)

func TestRowsToJson(t *testing.T) {
	tt := []struct {
		name        string
		mockFunc    func() *sql.Rows
		expected    []byte
		expectedErr error
	}{
		{
			name: "Success",
			mockFunc: func() *sql.Rows {
				db, mock, _ := sqlmock.New()
				defer db.Close()
				rows := sqlmock.NewRows([]string{"Column 1", "Column 2"})
				rows.AddRow("Dummy_1", 1)
				rows.AddRow("Dummy_2", 2)
				rows.AddRow("Dummy_3", 3)
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
				rs, _ := db.Query("SELECT 1")
				return rs
			},
			expected: []byte(`[{"Column_1":"Dummy_1","Column_2":1},{"Column_1":"Dummy_2","Column_2":2},{"Column_1":"Dummy_3","Column_2":3}]`),
		},
		{
			name: "Failure - Columns error",
			mockFunc: func() *sql.Rows {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("SELECT 1").WillReturnError(errors.New("rows is nil"))
				rs, _ := db.Query("SELECT 1")
				return rs
			},
			expectedErr: errors.New("rows is nil"),
		},
		{
			name: "Success - no rows",
			mockFunc: func() *sql.Rows {
				db, mock, _ := sqlmock.New()
				rows := sqlmock.NewRows([]string{"Col1", "Col2"})
				mock.ExpectQuery("SELECT 1").WillReturnRows(rows)
				rs, _ := db.Query("SELECT 1")
				return rs
			},
			expected: []byte(`[]`),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			bytes, err := sql2json.RowsToJson(tc.mockFunc())
			if err != nil && err.Error() != tc.expectedErr.Error() {
				t.Errorf("Expected error %v, but got %v", tc.expectedErr, err)
			}
			if string(bytes) != string(tc.expected) {
				t.Errorf("Expected %v, but got %v", string(tc.expected), string(bytes))
			}
		})
	}
}

func BenchmarkRowsToJson(b *testing.B) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	rows := sqlmock.NewRows([]string{"Column 1", "Column 2"})
	for i := 0; i < 1; i++ {
		rows.AddRow(fmt.Sprintf("Dummy_%d", i), i)
	}
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	rs, _ := db.Query("SELECT 1")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = sql2json.RowsToJson(rs)
	}

	b.StopTimer()
}

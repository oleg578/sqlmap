package sql2json

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"testing"
)

var (
	db                     *sql.DB
	con                    *sql.Conn
	stmt                   *sql.Stmt
	errDB, errCon, errStmt error
)

func prepareRows() (rows *sql.Rows, err error) {
	db, errDB = sql.Open("mysql", "root:admin@tcp(127.0.0.1:3307)/test")
	if errDB != nil {
		log.Fatal(errDB)
	}
	q := "SELECT * FROM `dummy` limit ?"
	ctx := context.Background()
	con, errCon = db.Conn(ctx)
	if errCon != nil {
		log.Fatal(errCon)
	}
	stmt, errStmt = con.PrepareContext(ctx, q)
	if errStmt != nil {
		log.Fatalf("stmt :%v", errStmt)
	}
	rs, errRS := stmt.QueryContext(ctx, 3)
	return rs, errRS
}

func BenchmarkRowsToJson(b *testing.B) {
	rs, errRs := prepareRows()
	defer db.Close()
	defer con.Close()
	defer stmt.Close()
	if errRs != nil {
		b.Fatal(errRs)
	}
	defer rs.Close()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := RowsToJson(rs)
		if err != nil {
			b.Fatalf("an error '%s' occurred during RowsToJson execution", err)
		}
	}
}

func mockRows() *sql.Rows {
	rs, errRs := prepareRows()
	if errRs != nil {
		log.Fatal(errRs)
	}
	return rs
}

// TestRowsToJson mainly tests the RowsToJson function with different cases.
func TestRowsToJson(t *testing.T) {
	cases := []struct {
		name          string
		mockRows      func() *sql.Rows
		expectedError error
	}{
		{
			name:          "Valid Rows",
			mockRows:      mockRows,
			expectedError: nil,
		},
	}
	for _, tc := range cases {
		rs, err := prepareRows()
		if err != nil {
			log.Fatal(err)
		}
		t.Run(tc.name, func(t *testing.T) {
			_, err := RowsToJson(rs)
			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("Test %s failed: expected error %v, got %v", tc.name, tc.expectedError, err)
			}
			if err == nil && tc.expectedError != nil {
				t.Errorf("Test %s failed: expected error %v, got nil", tc.name, tc.expectedError)
			}
		})
	}
}

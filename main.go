package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"runtime"
	"time"
)

func main() {
	Process()
}

func Process() {
	fmt.Println("process ...")
	db, mock, _ := sqlmock.New()
	defer db.Close()
	rows := sqlmock.NewRows([]string{"Column 1", "Column 2"})
	for i := 0; i < 10_000_000; i++ {
		rows.AddRow(fmt.Sprintf("Dummy_%d", i), i)
	}
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	rs, _ := db.Query("SELECT 1")
	// start time
	startTime := time.Now()
	// process rows
	_, err := RowsToJson(rs)
	if err != nil {
		panic(err)
	}
	//end time
	endTime := time.Now()
	fmt.Printf("Execution Time = %v\n", endTime.Sub(startTime).Milliseconds())
	printMemUsage()
}

type Dummy struct {
	Column1 string
	Column2 string
}

func RowsToJson(rows *sql.Rows) ([]byte, error) {
	var result = make([]Dummy, 0)
	if rows == nil {
		return nil, fmt.Errorf("rows is nil")
	}
	for rows.Next() {
		var rec = Dummy{}
		if err := rows.Scan(&rec.Column1, &rec.Column2); err != nil {
			return nil, err
		}
		result = append(result, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return json.Marshal(result)
}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
	fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	fmt.Printf("\tHeapAlloc = %v MiB", m.HeapAlloc/1024/1024)
	fmt.Printf("\tStackInuse = %v Kb", m.StackInuse/1024)
	fmt.Printf("\tSys = %v Kb", m.Sys/1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

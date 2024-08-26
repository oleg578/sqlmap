package main

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brianvoe/gofakeit/v7"
	"runtime"
	"sql2json"
	"time"
)

func main() {
	db, mock, _ := sqlmock.New()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	rows := sqlmock.NewRows([]string{"Id", "Product", "Price", "Qty", "NullData", "Date"})
	for i := 0; i < 10; i++ {
		rows.AddRow(
			i,
			gofakeit.Product().Name,
			gofakeit.Product().Price,
			gofakeit.IntRange(10, 1000),
			nil,
			gofakeit.Date().String())
	}
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	rs, _ := db.Query("SELECT 1")
	startTime := time.Now()
	msg, err := sql2json.RowsToJson(rs)
	if err != nil {
		panic(err)
	}
	endTime := time.Now()
	fmt.Printf("Elapsed time: %v ms\n", endTime.Sub(startTime).Milliseconds())
	printMemUsage()
	fmt.Println(len(msg))
	//fmt.Println(string(msg))
}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("TotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	fmt.Printf("\tHeapAlloc = %v MiB", m.HeapAlloc/1024/1024)
	fmt.Printf("\tStackInuse = %v Kb", m.StackInuse/1024)
	fmt.Printf("\tSys = %v Kb", m.Sys/1024)
	fmt.Printf("\tFrees = %v", m.Frees)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

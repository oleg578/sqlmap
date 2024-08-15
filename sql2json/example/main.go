package main

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"runtime"
	"sql2json"
	"time"
)

func main() {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	rows := sqlmock.NewRows([]string{"Column 1", "Column 2"})
	for i := 0; i < 10_000_000; i++ {
		rows.AddRow(fmt.Sprintf("Dummy_%d", i), i)
	}
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	rs, _ := db.Query("SELECT 1")
	startTime := time.Now()
	_, err := sql2json.RowsToJson(rs)
	if err != nil {
		panic(err)
	}
	endTime := time.Now()
	fmt.Printf("Elapsed time: %s\n", endTime.Sub(startTime))
	printMemUsage()

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

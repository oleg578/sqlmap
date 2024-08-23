package main

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-faker/faker/v4"
	"github.com/oleg578/sqlmap/sql2json"
	"runtime"
	"time"
)

func main() {
	const MOD = 10000007
	db, mock, _ := sqlmock.New()
	defer db.Close()
	rows := sqlmock.NewRows([]string{"Id", "Name", "Salary", "Hours", "Date"})
	for i := 0; i < 10_000_000; i++ {
		rows.AddRow(i, faker.FirstName(), 1000.00*i+(i*678)%MOD, i*8%MOD, faker.Timestamp())
	}
	//add null data
	//rows.AddRow(10, nil, nil, nil, nil)
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

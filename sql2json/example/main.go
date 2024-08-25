package main

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-faker/faker/v4"
	"log"
	"runtime"
	"sql2json"
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
	c, errCh := sql2json.RowsToJson(rs)

	for {
		select {
		case r, ok := <-c:
			if !ok {
				c = nil
			} else {
				_ = r
			}
		case err, ok := <-errCh:
			if ok {
				log.Printf("Error: %s", err)
			} else {
				errCh = nil
			}
		}

		// Exit the loop when both channels are closed
		if c == nil && errCh == nil {
			break
		}
	}

	endTime := time.Now()
	fmt.Printf("Elapsed time: %v ms\n", endTime.Sub(startTime).Milliseconds())
	printMemUsage()
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

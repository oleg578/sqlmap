package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"runtime"
	"sql2json"
	"time"
)

func main() {
	db, err := sql.Open("mysql", "root:admin@tcp(127.0.0.1:3307)/test")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	q := "SELECT * FROM `dummy`"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	con, errCon := db.Conn(ctx)
	if errCon != nil {
		log.Fatal(errCon)
	}
	defer con.Close()
	rs, err := con.QueryContext(ctx, q)
	if err != nil {
		log.Fatal(err)
	}
	defer rs.Close()
	fmt.Println("start parse rows...")
	startTime := time.Now()
	msg, errRTJ := sql2json.RowsToJson(rs)
	if errRTJ != nil {
		panic(errRTJ)
	}
	endTime := time.Now()
	fmt.Printf("Elapsed time: %v ms\n", endTime.Sub(startTime).Milliseconds())
	printMemUsage()
	fmt.Println(len(msg))
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

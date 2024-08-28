package main

import (
	"context"
	"database/sql"
	"github.com/brianvoe/gofakeit/v7"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	db, err := sql.Open("mysql", "root:admin@tcp(127.0.0.1:3307)/test")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	q := "INSERT INTO `dummy` VALUES(?,?,?,?,?,?)"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	con, err := db.Conn(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer con.Close()
	tx, errTx := con.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if errTx != nil {
		log.Fatal(errTx)
	}
	defer tx.Rollback()
	stmt, errStmt := tx.Prepare(q)
	if errStmt != nil {
		log.Fatal(errStmt)
	}
	for i := 0; i < 10_000_000; i++ {
		if _, err := stmt.ExecContext(ctx,
			i,
			gofakeit.Product().Name,
			gofakeit.Product().Price,
			gofakeit.IntRange(10, 1000),
			nil,
			gofakeit.Date().String()); err != nil {
			log.Fatal(err)
		}
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}

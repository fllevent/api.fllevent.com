package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func dbconnect(host string, database string, user string, password string) (db *sql.DB, err error) {
	db, err = sql.Open("mysql", user+":"+password+"@tcp("+host+":3306)/"+database+"?charset=utf8mb4,utf8&parseTime=True")
	return db, err
}

func dbclose(db *sql.DB) {
	db.Close()
	fmt.Println("database closed")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

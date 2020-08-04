package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var DB *sql.DB

func init() {
	DB, _ = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/f-file?charset=utf-8")
	DB.SetConnMaxLifetime(10)
	if err := DB.Ping(); err != nil {
		fmt.Println("open database fail")
		os.Exit(1)
	}
}

func DBConn() *sql.DB {
	return DB
}

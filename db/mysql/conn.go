package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var DB *sql.DB

func init() {
	DB, _ = sql.Open("mysql", "root:libo1121X@tcp(127.0.0.1:3306)/f-file?charset=utf8mb4&parseTime=true&loc=Local")
	DB.SetConnMaxLifetime(10)
	if err := DB.Ping(); err != nil {
		fmt.Println("open database fail", err)
		os.Exit(1)
	}
}

func DBConn() *sql.DB {
	return DB
}

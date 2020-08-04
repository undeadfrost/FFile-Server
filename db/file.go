package db

import (
	"FFile-Server/db/mysql"
	"fmt"
	"time"
)

func OnFileUploadFinished(fileSha1 string, fileName string, fileSize int64, fileAddr string) bool {
	stmt, err := mysql.DBConn().Prepare("insert ignore into file-list (`file_sha1`," +
		" `file_name`, `file_size`, `file_addr`, `status`, `create_at`, `update_at`) values (?, ?, ?, ?, 1, ?, ?)")

	if err != nil {
		fmt.Println("Failed to prepare statement, err: " + err.Error())
		return false
	}

	defer stmt.Close()

	createAt := time.Now().Unix()
	ret, err := stmt.Exec(fileSha1, fileName, fileSize, fileAddr, createAt, createAt)
	if err != nil {
		fmt.Printf(err.Error())
		return false
	}

	if rf, err := ret.RowsAffected(); err == nil {
		if rf < 0 {
			fmt.Printf("File with hash:%s has be upload before", fileName)
			return true
		}
	}

	return false
}

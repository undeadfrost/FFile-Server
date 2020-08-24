package db

import (
	db "FFile-Server/db/mysql"
	"fmt"
)

func OnUserUploadFinished(username, fileSha1, fileName string, fileSize int64) bool {
	stmt, err := db.DBConn().Prepare("insert ignore into `user_file` (" +
		"`username`, `file_sha1`, `file_name`, `file_size`) values (?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, fileSha1, fileName, fileSize)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if rf, err := ret.RowsAffected(); err == nil && rf > 0 {
		return true
	}

	fmt.Println(err.Error())
	return false
}

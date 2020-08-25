package db

import (
	db "FFile-Server/db/mysql"
	"fmt"
)

type UserFile struct {
	FileHash   string
	FileName   string
	FileSize   int64
	UploadAt   string
	LastUpdate string
}

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

func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	stmt, err := db.DBConn().Prepare("select file_sha1, file_name, file_size, upload_at, last_update" +
		" from `user_file` where username = ? limit ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, limit)
	if err != nil {
		return nil, err
	}

	var userFiles []UserFile
	for rows.Next() {
		userFile := UserFile{}
		err = rows.Scan(&userFile.FileHash, &userFile.FileName, &userFile.FileSize, &userFile.UploadAt, &userFile.LastUpdate)
		if err != nil {
			return nil, err
		}
		userFiles = append(userFiles, userFile)
	}

	return userFiles, nil
}

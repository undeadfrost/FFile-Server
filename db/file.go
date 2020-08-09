package db

import (
	"FFile-Server/db/mysql"
	"database/sql"
	"fmt"
	"time"
)

type TableFile struct {
	FileHash     string
	FileName     string
	FileSize     sql.NullInt64
	FileAddr     sql.NullString
	FileUpdateAt sql.NullTime
}

func OnFileUploadFinished(fileSha1 string, fileName string, fileSize int64, fileAddr string) bool {
	stmt, err := mysql.DBConn().Prepare("insert ignore into `file_list` (`file_sha1`," +
		" `file_name`, `file_size`, `file_addr`, `status`, `create_at`, `update_at`) values (?, ?, ?, ?, 1, ?, ?)")

	if err != nil {
		fmt.Println("Failed to prepare statement, err: " + err.Error())
		return false
	}

	defer stmt.Close()

	createAt := time.Now()
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

func GetFileMeta(fileHash string) (*TableFile, error) {
	stmt, err := mysql.DBConn().Prepare("select file_sha1, file_name, file_addr, file_size, update_at " +
		"from file_list where file_sha1 = ? and status = 1")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	var tableFile = TableFile{}
	err = stmt.QueryRow(fileHash).Scan(&tableFile.FileHash, &tableFile.FileName,
		&tableFile.FileAddr, &tableFile.FileSize, &tableFile.FileUpdateAt)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &tableFile, nil
}

func GetFileMetaList() []TableFile {
	stmt, err := mysql.DBConn().Prepare("select file_sha1, file_name, file_size, " +
		"file_addr, update_at from file_list")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	cols, err := rows.Columns()

	var tableFiles = make([]TableFile, 0, len(cols))
	for rows.Next() {
		var tableFile = TableFile{}

		err = rows.Scan(&tableFile.FileHash, &tableFile.FileName, &tableFile.FileSize, &tableFile.FileAddr, &tableFile.FileUpdateAt)
		if err != nil {
			fmt.Println(err.Error())
		}

		tableFiles = append(tableFiles, tableFile)
	}

	return tableFiles
}

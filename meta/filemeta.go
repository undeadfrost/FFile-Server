package meta

import (
	"FFile-Server/db"
)

type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

func UpdateFileMeta(f FileMeta) {
	fileMetas[f.FileSha1] = f
}

func UploadFileMetaDB(f FileMeta) {
	db.OnFileUploadFinished(f.FileSha1, f.FileName, f.FileSize, f.Location)
}
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

func RemoveFIleMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}

func GetFiles() map[string]FileMeta {
	return fileMetas
}

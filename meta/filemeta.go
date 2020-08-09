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

func UploadFileMetaDB(f FileMeta) bool {
	return db.OnFileUploadFinished(f.FileSha1, f.FileName, f.FileSize, f.Location)
}
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

func GetFileMetaDB(fileSha1 string) (FileMeta, error) {
	tableFile, err := db.GetFileMeta(fileSha1)
	if err != nil {
		return FileMeta{}, err
	}

	fileMeta := FileMeta{
		FileSha1: tableFile.FileHash,
		FileName: tableFile.FileName,
		FileSize: tableFile.FileSize.Int64,
		Location: tableFile.FileAddr.String,
		UploadAt: tableFile.FileUpdateAt.Time.Format("2006-01-02 15:04:05"),
	}

	return fileMeta, nil
}

func RemoveFIleMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}

func GetFiles() map[string]FileMeta {
	return fileMetas
}

func GetFilesDB() []FileMeta {
	tableFiles := db.GetFileMetaList()
	var fileMetas = make([]FileMeta, 0, len(tableFiles))

	for k, _ := range tableFiles {
		fileMeta := FileMeta{
			FileSha1: tableFiles[k].FileHash,
			FileName: tableFiles[k].FileName,
			FileSize: tableFiles[k].FileSize.Int64,
			Location: tableFiles[k].FileAddr.String,
			UploadAt: tableFiles[k].FileUpdateAt.Time.Format("2006-01-02 15:04:05"),
		}
		fileMetas = append(fileMetas, fileMeta)
	}

	return fileMetas
}

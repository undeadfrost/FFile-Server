package handler

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"os"
	"time"

	"FFile-Server/meta"
	"FFile-Server/util"
)

func UploadHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	file, head, err := r.FormFile("file")
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		return
	}
	defer file.Close()

	fileMeta := meta.FileMeta{
		FileName: head.Filename,
		Location: "./upload/" + head.Filename,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	destFile, err := os.Create(fileMeta.Location)
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		return
	}
	defer destFile.Close()

	fileMeta.FileSize, err = io.Copy(destFile, file)
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		return
	}

	fileMeta.FileSha1 = util.FileSha1(destFile)
	meta.UpdateFileMeta(fileMeta)

	io.WriteString(w, "Upload finished!")
}

func GetFileMetaHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fileHash := ps.ByName("fileHash")
	fileMeta := meta.GetFileMeta(fileHash)

	data, err := json.Marshal(fileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Printf("t")
	fmt.Fprintf(w, "hello world")
}

func UserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello %s", ps.ByName("name"))
}

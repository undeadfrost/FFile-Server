package handler

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"FFile-Server/meta"
	"FFile-Server/util"
)

type updateBody struct {
	FileName string
}

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
	_ = meta.UploadFileMetaDB(fileMeta)

	io.WriteString(w, "Upload finished!")
}

func GetFilesHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	files := meta.GetFilesDB()
	data, err := json.Marshal(files)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func GetFileMetaHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fileHash := ps.ByName("fileHash")
	fileMeta, err := meta.GetFileMetaDB(fileHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(fileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func DownloadFileHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fileHash := ps.ByName("fileHash")
	fileMeta := meta.GetFileMeta(fileHash)

	file, err := os.Open(fileMeta.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileName := url.QueryEscape(fileMeta.FileName)
	w.Header().Set("Content-type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	w.Write(data)
}

func DeleteFileHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fileHash := ps.ByName("fileHash")
	fileMeta := meta.GetFileMeta(fileHash)

	os.Remove(fileMeta.Location)
	meta.RemoveFIleMeta(fileHash)

	w.WriteHeader(http.StatusOK)
}

func UpdateFileHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var updateData = new(updateBody)
	fileHash := ps.ByName("fileHash")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &updateData)
	if err != nil {
		fmt.Printf("Unmarshal err, %v\n", err)
		return
	}

	curFileMeta := meta.GetFileMeta(fileHash)
	curFileMeta.FileName = updateData.FileName
	meta.UpdateFileMeta(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Printf("t")
	fmt.Fprintf(w, "hello world")
}

func UserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello %s", ps.ByName("name"))
}

package handler

import (
	"FFile-Server/db"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"FFile-Server/meta"
	"FFile-Server/util"
)

type updateBody struct {
	FileName string
}

type tryFastBody struct {
	FileHash string `json:"fileHash"`
	FileName string `json:"fileName"`
}

func UploadHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 1.创建文件上传目录
	folder := "./upload/" + time.Now().Format("20060102")
	isFolder, _ := util.PathExists(folder)
	if !isFolder {
		if err := os.MkdirAll(folder, 777); err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// 2.创建文件,保存文件不带后缀名
	fn := folder + "/" + uuid.NewV4().String()
	fnTmp, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 777)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	file, head, err := r.FormFile("file")
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		return
	}
	defer file.Close()

	fileSize, err := io.Copy(fnTmp, file)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileMeta := meta.FileMeta{
		FileName: head.Filename,
		FileSize: fileSize,
		Location: fn,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	// 3.上传文件已存在则删除本次上传
	fp, _ := os.Open(fn)
	fileMeta.FileSha1 = util.FileSha1(fp)
	uploadSuc, err := db.OnFileUploadFinished(fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize, fileMeta.Location)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer func() {
		_ = fp.Close()
		_ = fnTmp.Close()
		if !uploadSuc {
			fmt.Println("Delete File", fn)
			err = os.Remove(fn)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}()

	// 4.保存文件到用户表
	username := r.Context().Value("username").(string)
	suc := db.OnUserUploadFinished(username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize)
	if suc {
		io.WriteString(w, "Upload Finished!")
	} else {
		io.WriteString(w, "Upload Failed")
	}
}

func TryFastUploadHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 1.获取参数
	var tryFastBody = tryFastBody{}
	username := r.Context().Value("username").(string)
	err := json.NewDecoder(r.Body).Decode(&tryFastBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 2.查询文件是否已经存在
	tableFile, err := meta.GetFileMetaDB(tryFastBody.FileHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 3.文件不存在无法秒传
	if tableFile == nil {
		w.Header().Set("Content-type", "application/json")
		rawRep := util.AjaxReturn(0, "success", nil)
		w.Write(rawRep.JsonBytes())
		return
	}

	// 4.查询用户表是否存在相同文件
	isExist, err := db.FileExist(username, tryFastBody.FileHash, tryFastBody.FileName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if isExist {
		fileExt := path.Ext(tryFastBody.FileName)
		index := strings.LastIndex(tryFastBody.FileName, fileExt)
		timeNow := time.Now().Unix()
		tryFastBody.FileName = tryFastBody.FileName[:index] + strconv.FormatInt(timeNow, 10) + tryFastBody.FileName[index:]
	}
	// 5.写入用户表
	suc := db.OnUserUploadFinished(username, tryFastBody.FileHash, tryFastBody.FileName, tableFile.FileSize)

	w.Header().Set("Content-type", "application/json")

	if suc {
		rawRep := util.AjaxReturn(0, "success", nil)
		w.Write(rawRep.JsonBytes())
	} else {
		rawRep := util.AjaxReturn(-1, "Failed", nil)
		w.Write(rawRep.JsonBytes())
	}
}

func GetFilesHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username := r.Context().Value("username").(string)
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userFiles, err := db.QueryUserFileMetas(username, limit)
	// files := meta.GetFilesDB()
	// data, err := json.Marshal(files)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rawRep := util.AjaxReturn(0, "success", userFiles)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(rawRep.JsonBytes())
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

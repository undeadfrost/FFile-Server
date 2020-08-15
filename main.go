package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"FFile-Server/handler"
)

func main() {
	router := httprouter.New()

	router.GET("/", handler.IndexHandler)
	router.GET("/user/:name", handler.UserHandler)
	router.POST("/upload", handler.UploadHandler)
	router.GET("/file", handler.GetFilesHandler)
	router.GET("/file/:fileHash", handler.GetFileMetaHandler)
	router.GET("/download/:fileHash", handler.DownloadFileHandler)
	router.DELETE("/file/:fileHash", handler.DeleteFileHandler)
	router.PUT("/file/:fileHash", handler.UpdateFileHandler)

	router.POST("/user/signin", handler.SignIn)
	router.POST("/user/login", handler.Login)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("Failed to start server error: $s", err.Error())
	}
}

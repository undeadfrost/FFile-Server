package main

import (
	"FFile-Server/handler"
	"FFile-Server/middleware"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func main() {
	router := httprouter.New()

	router.GET("/", handler.IndexHandler)
	router.GET("/welcome/:name", handler.UserHandler)
	router.POST("/upload", middleware.CheckLogin(handler.UploadHandler))
	router.GET("/file", handler.GetFilesHandler)
	router.GET("/file/:fileHash", handler.GetFileMetaHandler)
	router.GET("/download/:fileHash", handler.DownloadFileHandler)
	router.DELETE("/file/:fileHash", handler.DeleteFileHandler)
	router.PUT("/file/:fileHash", handler.UpdateFileHandler)

	router.POST("/user/signin", handler.SignIn)
	router.POST("/user/login", handler.Login)
	router.GET("/user/info", middleware.CheckLogin(handler.UserInfo))

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("Failed to start server error: $s", err.Error())
	}
}

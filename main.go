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
	router.POST("/file", handler.UploadHandler)
	router.GET("/file/:fileHash", handler.GetFileMetaHandler)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("Failed to start server error: $s", err.Error())
	}
}

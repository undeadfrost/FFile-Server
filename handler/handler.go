package handler

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"os"
)

func UploadHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	file, head, err := r.FormFile("file")
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		return
	}
	defer file.Close()

	destFile, err := os.Create("./upload/" + head.Filename)
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		return
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, file)
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		return
	}

	io.WriteString(w, "Upload finished!")
}

func IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Printf("t")
	fmt.Fprintf(w, "hello world")
}

func UserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello %s", ps.ByName("name"))
}

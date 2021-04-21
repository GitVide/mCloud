package main

import (
	"Filestore-SERVER/handler"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc",handler.UploadSucHandler)
	err:=http.ListenAndServe(":8080",nil)
	if err != nil {
		fmt.Printf("Fail to start server,err:%s", err.Error())
	}
}


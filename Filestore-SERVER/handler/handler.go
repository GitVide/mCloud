package handler

import (
	"Filestore-SERVER/meta"
	"Filestore-SERVER/util"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// UploadHandler 处理文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//返回上传html页面
		data,err := ioutil.ReadFile("./static/view/index.html")
		if err != nil{
			io.WriteString(w,"internal server error")
			return
		}
		io.WriteString(w,string(data))
	}else if r.Method == "POST" {
		//接收文件流及存储到本地目录
		file, head,err:=r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get data,err:%s\n", err.Error())
			return
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "/tmp/"+head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		newFile,err:=os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("Failed to create file,err:%s\n",err.Error())
		}
		defer newFile.Close()

		fileMeta.FileSize,err=io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to save data into file,err:%s\n",err.Error())
			return
		}

		newFile.Seek(0,0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		meta.UpdateFileMeta(fileMeta)

		http.Redirect(w,r,"/file/upload/suc",http.StatusFound)
	}
}

// UploadSucHandler 上传已完成
func UploadSucHandler(w http.ResponseWriter, r*http.Request){
	io.WriteString(w,"Upload finished")
}

// GetFileMetaHandler 获取文件元信息接口
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request){

	r.ParseForm()
	filehash:=r.Form["filehash"][0]
	fMeta:=meta.GetFileMeta(filehash)
	data,err:=json.Marshal(fMeta)
	if err!=nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// DownloadHandler 下载文件接口
func DownloadHandler(w http.ResponseWriter,r *http.Request)  {
	r.ParseForm()
	fsha1:=r.Form.Get("filehash")
	fm:=meta.GetFileMeta(fsha1)

	f,err:=os.Open(fm.Location)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data,err:=ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type","application/octect-stream")
	w.Header().Set("Content-disposition","attachment;filename=\""+fm.FileName+"\"")
	w.Write(data)
}

// FileMetaUpdateHandler 文件重命名接口
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	opType:=r.Form.Get("op")
	fileSha1:=r.Form.Get("filehash")
	newFileName:=r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method != "POST"{
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta:=meta.GetFileMeta(fileSha1)
	curFileMeta.FileName=newFileName
	meta.UpdateFileMeta(curFileMeta)

	w.WriteHeader(http.StatusOK)
	data,err:=json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// FileDeleteHandler 文件删除接口
func FileDeleteHandler(w http.ResponseWriter,r*http.Request)  {
	r.ParseForm()
	fileSha1:=r.Form.Get("filehash")
	fMeta:=meta.GetFileMeta(fileSha1)
	os.Remove(fMeta.Location)
	meta.RemoveFileMeta(fileSha1)
	w.WriteHeader(http.StatusOK)
}
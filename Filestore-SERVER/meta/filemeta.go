package meta

// FileMeta 文件元信息结构
type FileMeta struct {
	//文件的唯一标志
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas=make(map[string]FileMeta)
}

// UpdateFileMeta 新增/更新文件元信息
func UpdateFileMeta(fMeta FileMeta) {
	fileMetas[fMeta.FileSha1] = fMeta
}

// GetFileMeta 通过fileSha1获取文件元信息
func GetFileMeta(fileSha1 string) FileMeta{
	return fileMetas[fileSha1]
}

// RemoveFileMeta 删除元信息
func RemoveFileMeta(fileSha1 string)  {
	delete(fileMetas,fileSha1)
}

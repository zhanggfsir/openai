package commonUtil

import (
	"github.com/astaxie/beego/logs"
	"github.com/satori/go.uuid"
	"os"
)

func GenUuid() string {
	// 创建可以进行错误处理的 UUID v4
	uuidRes := uuid.NewV4()
	return uuidRes.String()
}

func FileExt(path string) string {
	for i := len(path) - 1; i >= 0 && path[i] != '/'; i-- {
		if path[i] == '.' {
			return path[i:]
		}
	}
	return ""
}

//获得文件保存的全路径:dir/fileName
func CreateFullPath(dir string) error {
	//1首先创建保存文件的文件夹
	if err := os.MkdirAll(dir , os.ModePerm); err != nil {
		logs.Error("创建文件夹错误err= ", err)
		//logs.Error("创建文件夹错误err= ", err)
		return err
	}
	return  nil
}
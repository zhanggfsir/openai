package commonUtil

import (
	"github.com/astaxie/beego/logs"
	"os"
)

// xx/oo 如果oo存在就删除，然后继续创建
func ReCreateDir(saveFileDir string) error {
	err :=os.RemoveAll(saveFileDir)
	if err == nil {
		err:=os.MkdirAll(saveFileDir,0755)
		return err
	}
	return err
}

func CreateDir(dir string) error {
	err := os.MkdirAll(dir, 0755)
	logs.Info("生成路徑成功:",dir)
	return err
}

// xx/oo 如果oo不存在就删除，存在不作处理
func PathExists(saveFileDir string)  error {
	_, err := os.Stat(saveFileDir)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		err =os.Mkdir(saveFileDir,0755)
		if err!=nil{
			logs.Error("创建保存目录失败")
			return err
		}
	}
	return err
}
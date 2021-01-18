package commonUtil

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego/logs"
	"io"
	"os"
	"path/filepath"
)

func FileMd5(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		logs.Error("err= ", err)
	}

	md5 := md5.New()
	io.Copy(md5,file)
	return hex.EncodeToString(md5.Sum(nil))
}

func CopyFileToDir(srcFileDir,destDir string) error {
	srcFile, err := os.Open(srcFileDir)
	if err!=nil{
		logs.Info(err)
		return err
	}

	dstFile, err := os.OpenFile(destDir+"/"+filepath.Base(srcFile.Name()), os.O_WRONLY|os.O_CREATE, 0666)
	logs.Info(destDir+"/"+filepath.Base(srcFile.Name()))
	if err != nil {
		logs.Error(err)
		return err
	}

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		logs.Error(err)
		return err
	}
	srcFile.Close()
	dstFile.Close()

	return nil
}

//
func CopyFileFromDisk(srcDir string, fileName string, distDir string) error {
	srcFile, err := os.Open(srcDir + "/" + fileName)
	if err != nil {
		return err
	}

	//打开dstFileName
	dstFile, err := os.OpenFile(distDir+"/"+ fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	srcFile.Close()
	dstFile.Close()
	return nil
}




//func copyFile2Dir(dstFileName string, srcFileName string)   error {
//
//	srcFile, err := os.Open(srcFileName)
//	if err != nil {
//		return  err
//	}
//	defer srcFile.Close()
//
//	//通过srcFile，获取到READER
//	reader := bufio.NewReader(srcFile)
//
//	//打开dstFileName
//	dstFile, err := os.OpenFile(dstFileName+"/"+path.Base(srcFileName), os.O_WRONLY|os.O_CREATE, 0666)
//	if err != nil {
//		return err
//	}
//
//	//通过dstFile，获取到WRITER
//	writer := bufio.NewWriter(dstFile)
//	//writer.Flush()
//
//	defer dstFile.Close()
//
//	_,err =io.Copy(writer, reader)
//	writer.Flush()
//	return err
//}






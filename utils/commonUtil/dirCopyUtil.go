package commonUtil

import (
	"bufio"
	"openai-backend/utils/stringUtil"
	"golang.org/x/sync/errgroup"
	"io"
	"io/ioutil"
	"os"
	"path"
)

/**
 * 拷贝文件夹中的文件 xx/A/ --> xx/B/
 * @param srcPath  		需要拷贝的文件夹路径: D:/test
 * @param destPath		拷贝到的位置: D:/backup/
 */
func CopyDir(destPath,srcPath string) <- chan  error {
	returnErr :=make(chan error ,1)
	var defaultZipTypes = []string{".json",".bmp", ".jpg", ".png"} // gz rar tar  // "zip","rar","7z","tar","gz","bz2"

	go func() {
		fileInfo, err := ioutil.ReadDir(srcPath)
		var eg errgroup.Group
		//eg.Add(len(fileInfo))
		for i := range fileInfo {
			//eg.Go()
			go func(file os.FileInfo) {
				eg.Go(func() error {  //
						//wg.Done()
						index := stringUtil.StringsContains(defaultZipTypes, path.Ext(file.Name()))  // 找不到返回 -1  找到返回数组下标
						if !file.IsDir()  && index>-1 {
							err = copyFile2Dir(destPath, srcPath + "/" + file.Name())
							if err != nil {
								return  err  // 有一个报错就会中断
							}
						}
						return nil
					})
			}(fileInfo[i])
		}
		returnErr <- eg.Wait()  // err收集 nil
	}()


	return returnErr
}

func CopyDir2(destPath,srcPath string) error{
	//returnErr :=make(chan error ,1)
	var defaultZipTypes = []string{".json",".bmp", ".jpg", ".png"} // gz rar tar  // "zip","rar","7z","tar","gz","bz2"

	//go func() {
		fileInfo, err := ioutil.ReadDir(srcPath)
		for _, file := range fileInfo {
			index := stringUtil.StringsContains(defaultZipTypes, path.Ext(file.Name()))  // 找不到返回 -1  找到返回数组下标
			if !file.IsDir()  && index>-1 {
				err = copyFile2Dir(destPath, srcPath + "/" + file.Name())
				if err != nil {
					return err
				}
			}
		}
	//}()

	return nil
}

func CopyDir3(destPath,srcPath string) <- chan  error {
	returnErr :=make(chan error ,1)
	var defaultZipTypes = []string{".json",".bmp", ".jpg", ".png"} // gz rar tar  // "zip","rar","7z","tar","gz","bz2"

	go func() {
		fileInfo, err := ioutil.ReadDir(srcPath)
		for _, file := range fileInfo {
			index := stringUtil.StringsContains(defaultZipTypes, path.Ext(file.Name()))  // 找不到返回 -1  找到返回数组下标
			if !file.IsDir()  && index>-1 {
				err = copyFile2Dir(destPath, srcPath + "/" + file.Name())
				if err != nil {
					returnErr <- err
					//return returnErr
				}
			}
		}
	}()

	return returnErr
}

func copyFile2Dir(dstFileName string, srcFileName string)   error {

	srcFile, err := os.Open(srcFileName)
	if err != nil {
		return  err
	}
	defer srcFile.Close()

	//通过srcFile，获取到READER
	reader := bufio.NewReader(srcFile)

	//打开dstFileName
	dstFile, err := os.OpenFile(dstFileName+"/"+path.Base(srcFileName), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	//通过dstFile，获取到WRITER
	writer := bufio.NewWriter(dstFile)
	//writer.Flush()

	defer dstFile.Close()

	_,err =io.Copy(writer, reader)
	writer.Flush()
	return err
}

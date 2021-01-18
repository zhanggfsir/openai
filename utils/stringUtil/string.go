package stringUtil

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func TrimEnhance(body string) string {
	body = strings.Replace(body, " ", "", -1)
	body = strings.Replace(body, "\n", "", -1)
	body = strings.Replace(body, "\r", "", -1)
	body = strings.Replace(body, "\t", "", -1)
	return body
}

//suffix.jpg  ->jpg
func SuffixFileName(path string) (dir string) {
	i := strings.LastIndex(path, ".")
	if i < 0 {
		return path
	}
	return path[i+1:]
}

//prefix.jpg  ->prefix 		 pic/prefix.jpg  ->pic/prefix
func PrefixFileName(path string) (dir string) {
	i := strings.LastIndex(path, ".")
	if i < 0 {
		return path
	}
	return path[:i]
}

//pic2/x.jpg ->x.jpg
func FileNameSplit(path string) (dir string) {
	i := strings.LastIndex(path, "/")
	if i < 0 {
		return path
	}
	return path[i+1:]
}
//pic/x.jpg ->pic/
func PathSplit(path string) (dir string) {
	i := strings.LastIndex(path, "/")
	if i < 0 {
		return path
	}
	return path[:i+1]
}

//pic/prefix.jpg  ->pic/prefix
func FullPathSplit(path string) (dir string) {
	i := strings.LastIndex(path, ".")
	if i < 0 {
		return path
	}
	return path[:i]
}


//pic2/x.jpg -> x
func FileNamePrefixSplit(path string) (fileNamePrefix string) {
	if path==""{
		return ""
	}
	i := strings.LastIndex(path, "/")
	pureFileName:=path[i+1:]
	j:=strings.LastIndex(pureFileName,".")
	return  pureFileName[:j]
}

// jpg in {"bmp", "jpg", "png"}
func StringsContains(array []string, val string) (index int) {
	index = -1
	for i := 0; i < len(array); i++ {
		if array[i] == val {
			index = i
			return
		}
	}
	return
}


func Md5(str string) string  {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

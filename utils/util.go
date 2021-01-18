package utils

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/satori/go.uuid"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	NUmStr  = "0123456789"
	CharStr = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	SpecStr = "+=-@#~[]()!%^*$"
)

type Fields map[string]interface{}

var defaultR = rand.New(rand.NewSource(time.Now().Unix()))

func GenerateUuid() string {
	// 创建可以进行错误处理的 UUID v4
	uuidRes := uuid.NewV4()
	return uuidRes.String()
}

func MapIsnull(contentBody map[string]string, Key string) bool {
	if value, ok := contentBody[Key]; !ok || value == "" {
		return false
	}
	return true
}

func ToLower(str string) string {
	return strings.ToLower(str)
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func A2String(value interface{}) string {
	switch value.(type) {
	case string:
		return value.(string)
	case int64:
		var valueInt64 = value.(int64)
		return strconv.FormatInt(valueInt64, 10)
	case float64:
		var valueInt64 = value.(float64)
		return strconv.FormatFloat(valueInt64, 'f', -1, 64)
	default:
		logs.Info("以上类型都不是 type", value)
		return ""
	}
}

func A2Int64(value interface{}) (int64, error) {
	switch value.(type) {
	case string:
		valueInt64, err := strconv.ParseInt(value.(string), 10, 64)
		return valueInt64, err
	case int64:
		var valueInt64 = value.(int64)
		return valueInt64, nil
	default:
		logs.Error("以上类型都不是 type", value)
		return 0, errors.New("参数的类型错误")
	}
}

func A2Int(value interface{}) (int, error) {
	switch value.(type) {
	case string:
		valueInt, err := strconv.Atoi(value.(string))
		return valueInt, err
	case int:
		valueInt := value.(int)
		return valueInt, nil
	default:
		logs.Error("以上类型都不是 type", value)
		return 0, errors.New("参数的类型错误")
	}
}

func A2Float64(value interface{}) (float64, error) {
	switch value.(type) {
	case string:
		valueFloat64, err := strconv.ParseFloat(value.(string),64)
		if err != nil {
			logs.Error("数据类型转换错误err= ", err, "value= ", value)
			valueFloat64 = 0
		}
		return valueFloat64, err
	case float64:
		var valueFloat64 = value.(float64)
		return valueFloat64, nil
	default:
		logs.Error("以上类型都不是 value=", value)
		return 0, errors.New("参数的类型错误")
	}
}

//获得固定位数的字符串
func RandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := defaultR.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

//生成密码
//num:只使用数字[0-9],
//char:只使用英文字母[a-zA-Z]
//mix:使用数字和字母
//advance:使用数字、字母以及特殊字符`
func GenerateRandom(charset string, length int) string {
	//初始化密码切片
	var passwd []byte = make([]byte, length, length)
	var sourceStr string
	//判断字符类型,如果是数字
	if charset == "num" {
		sourceStr = NUmStr
		//如果选的是字符
	} else if charset == "char" {
		sourceStr = charset
		//如果选的是混合模式
	} else if charset == "mix" {
		sourceStr = fmt.Sprintf("%s%s", NUmStr, CharStr)
		//如果选的是高级模式
	} else if charset == "advance" {
		sourceStr = fmt.Sprintf("%s%s%s", NUmStr, CharStr, SpecStr)
	} else {
		sourceStr = NUmStr
	}
	//遍历，生成一个随机index索引,
	for i := 0; i < length; i++ {
		index := rand.Intn(len(sourceStr))
		passwd[i] = sourceStr[index]
	}
	return string(passwd)
}

func preNUm(data byte) int {
	str := fmt.Sprintf("%b", data)
	var i int = 0
	for i < len(str) {
		if str[i] != '1' {
			break
		}
		i++
	}
	return i
}

func IsUtf8(data []byte) bool {
	for i := 0; i < len(data);  {
		if data[i] & 0x80 == 0x00 {
			// 0XXX_XXXX
			i++
			continue
		} else if num := preNUm(data[i]); num > 2 {
			// 110X_XXXX 10XX_XXXX
			// 1110_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_0XXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_10XX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_110X 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// preNUm() 返回首个字节的8个bits中首个0bit前面1bit的个数，该数量也是该字符所使用的字节数
			i++
			for j := 0; j < num - 1; j++ {
				//判断后面的 num - 1 个字节是不是都是10开头
				if data[i] & 0xc0 != 0x80 {
					return false
				}
				i++
			}
		} else  {
			//其他情况说明不是utf-8
			return false
		}
	}
	return true
}

//获得时间戳
func GetTimestamp(secType string) int64 {
	if secType == "millSec" {
		return time.Now().UnixNano() / 1e6
	} else {
		return time.Now().Unix()
	}
}

//email verify
func VerifyEmailFormat(email string) bool {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

//mobile verify
func VerifyMobile(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|(19[0,3,5-8])|(147))\\d{8}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

//截取字符串 start 起点下标 end 终点下标(不包括)
func Substr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		return str
	}

	if end < 0 || end > length {
		return str
	}
	return string(rs[start:end])
}

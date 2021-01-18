// package conf 为配置相关.
package utils

import (
	"path/filepath"
)

//用户密码加密
const (
	PW_HASH_BYTES_LEN = 8
	PasswordSalt = "hdyusiIdSI^SOl"
)

const (
	JWT_KEY                         string = "yu8*uYSTJCGJHSM"
	JWT_DEFAULT_EXPIRE_SECONDS      int64    = 3600 //uint:sec default 60 minutes
	JWT_DEFAULT_LONG_EXPIRE_SECONDS int64    = 3600*24*1 // default 1 days
)

var (
	ConfigurationFile = "./conf/app.conf"
	WorkingDirectory  = "./"
	LogFile           = ""
	BaseUrl           = ""
	AutoLoadDelay     = 0
	DefaultDownloadUrlPre = "http://127.0.0.1:8091/"
)

func WorkingDir(elem ...string) string {
	elements := append([]string{WorkingDirectory}, elem...)
	return filepath.Join(elements...)
}

func init() {
	if p, err := filepath.Abs("./conf/app.conf"); err == nil {
		ConfigurationFile = p
	}
	if p, err := filepath.Abs("./"); err == nil {
		WorkingDirectory = p
	}
}

package utils

import (
	"encoding/base64"
	"fmt"
	"github.com/astaxie/beego/logs"
	"golang.org/x/crypto/scrypt"
)

/*
@Author: wangc293
@Time: 2020-05-30 10:02
@Description:
*/

func GenPassword(password string) string {
	salt := []byte(PasswordSalt)
	hash, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, PW_HASH_BYTES_LEN)
	if err != nil {
		logs.Error(err)
	}
	return fmt.Sprintf("%x", hash)
}

func Base64Encode(str string) string {
	encodeData := []byte(str)
	return base64.StdEncoding.EncodeToString(encodeData)
}
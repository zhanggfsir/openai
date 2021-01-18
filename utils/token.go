package utils

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/dgrijalva/jwt-go"
	"openai-backend/utils/httpUtil"
	"strings"
	"time"
)

/*
@Author:
@Time: 2020-07-30 10:02
@Description:
*/

type JwtUser struct {
	Id       string `json:"id"`
	Type     string `json:"type"` //inner  third
	Domain   string `json:"domain"`
	UserType string `json:"user_type"`
}

type MyCustomClaims struct {
	JwtUser
	jwt.StandardClaims
}

func ValidateToken(tokenString string) (JwtUser, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(JWT_KEY), nil
		})

	if token != nil{
		//logs.Info("TOKEN EXPIRE AT: ",time.Unix(token.Claims.(*MyCustomClaims).ExpiresAt,0).Format("2006-01-02 15:04:05"))
		//logs.Info(*token.Claims.(*MyCustomClaims))
	}else {
		logs.Info("token is nil")
	}

	//判断token验证的结果
	if err != nil {
		logs.Error("err=", err)
		value, ok := err.(*jwt.ValidationError)
		if ok && value.Errors == jwt.ValidationErrorExpired {
			return JwtUser{}, errors.New(httpUtil.ERROR_TOKEN_EXPIRED)
		}
		return JwtUser{}, errors.New(httpUtil.ERROR_TOKEN_INVALID)
	}
	//返回结果
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims.JwtUser, nil
	} else {
		logs.Error("validate tokenString failed !!!", err)
		return JwtUser{}, errors.New(httpUtil.ERROR_TOKEN_INVALID)
	}
}

func GenerateToken(user JwtUser, expiredSeconds int64) (tokenString string) {
	if expiredSeconds == 0 {
		expiredSeconds = JWT_DEFAULT_EXPIRE_SECONDS
	}

	mySigningKey := []byte(JWT_KEY)
	expireAt := time.Now().Add(time.Second * time.Duration(expiredSeconds)).Unix()


	logs.Info("TOKEN DURATION:",expiredSeconds,"TOKEN,EXPIRE AT: ",
		time.Now().Add(time.Second * time.Duration(expiredSeconds)).Format("2006-01-02 15:04:05"))

	claims := MyCustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: expireAt,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(mySigningKey)
	if err != nil {
		logs.Error("generate json web token failed !! error :", err)
	}
	return tokenStr
}

// 校验header中的token
func VerifySession(headerContent map[string][]string) (JwtUser, error) {
	if value, ok := headerContent["Authorization"]; ok {
		authorization := value[0]
		if strings.Contains(authorization, "Bearer") {
			if len(authorization) > 8 {
				token := string(authorization[7:])
				userData, err := ValidateToken(token)
				return userData, err
			} else {
				logs.Error("Token error=", authorization)
				return JwtUser{}, errors.New(httpUtil.ERROR_TOKEN_INVALID)
			}
		} else {
			logs.Error("Token error=", authorization)
			return JwtUser{}, errors.New(httpUtil.ERROR_TOKEN_NOT_EXIST)
		}
	} else {
		logs.Error("Token error=", headerContent)
		return JwtUser{}, errors.New(httpUtil.ERROR_TOKEN_NOT_EXIST)
	}
}

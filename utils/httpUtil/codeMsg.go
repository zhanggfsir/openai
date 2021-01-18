package httpUtil

const (
	//PUBLIC
	SUCCESS                     = "0"      //成功
	SYSTEM_ERROR                = "500"    //系统错误
	ERROR_TOKEN_INVALID         = "1000"   //token无效
	ERROR_TOKEN_EXPIRED         = "1001"   //token过期
	ERROR_TOKEN_NOT_EXIST       = "1002"   //token不存在
	ERROR_REFRESH_TOKEN_EXPIRED = "1003"   //刷新token过期
	ERROR_NO_AUTHORITY          = "1004"   //接口请求没有权限
	PARAMETER_ERROR             = "1005"   //接口请求参数错误
	ERROR_APPKEY_NOT_EXIST      = "1006"   //appKey传值为空
	ERROR_APPKEY_INVALID        = "1007"   //appKey无效
	ERROR_APPSECRET_INVALID     = "1008"   //appSecret无效
	ERROR_JSON_PATTERN          = "1009"   //json格式错误
	ERROR_CAPTCHA               = "1010"   //验证码错误
	ERROR_DOMAIN_NOT_EXIST      = "1010"   //域名不存在
	CREATE_DIR_ERROR = "1013"

	//用户登录
	ERROR_USER_LOGIN_ERROR      = "2001"
	ERROR_USER_NOT_EXIST        = "2002"

	OLD_PASSWORD_ERROR          = "2004"

	DEVICE_MAC_EMPTY_ERROR      = "2006"
	DEVICE_MAC_NO_FOUND_ERROR   = "2007"
	DEVICE_MAC_ERROR            = "2008"

	DEVICE_MAC_FORMAT_ERROR     = "2009"


	// import

	//export

	//  datasetTask
	InsertDatasetTaskError="7001"

	// dataInfo
	InsertDataInfoError="8001"

	// deal File
	WRITE_FILE_ERROR          = "3001"
	READ_FILE_ERROR           = "3002"
	GET_FILE_SIZE_ERROR       = "3003"
	OPEN_FILE_ERROR           = "3004"
	ILLEGAL_FILE_FORMAT       = "3005"
	SAVE_FILE_ERROR           = "3006"
	FILE_SIZE_TOO_LARGE		  = "3007"
	FILE_COMPRESS_ERROR		  = "3008"


	// zip
	CREATE_ZIP_ERROR		  ="3050"
	GET_FILE_HEADER_ERROR	  ="3051"
	WRITE_ZIP_ERROR			  ="3052"
)

var MsgFlags = map[string]string{
	SUCCESS:                          "success",
	SYSTEM_ERROR:                     "system error,please contact the admin",
	ERROR_TOKEN_INVALID:              "token invalid",
	ERROR_TOKEN_NOT_EXIST:            "token not exist",
	ERROR_TOKEN_EXPIRED:              "token expired",
	ERROR_REFRESH_TOKEN_EXPIRED:      "refresh token expired",
	ERROR_NO_AUTHORITY:               "没有权限",
	PARAMETER_ERROR:                  "parameter error",
	ERROR_APPKEY_NOT_EXIST:           "appKey not exist",
	ERROR_APPKEY_INVALID:             "appKey invalid",
	ERROR_APPSECRET_INVALID:          "appSecret invalid",
	ERROR_JSON_PATTERN:               "json error",
	ERROR_DOMAIN_NOT_EXIST:           "domain not exist",



	DEVICE_MAC_FORMAT_ERROR:          "MAC地址格式错误",
	DEVICE_MAC_EMPTY_ERROR:           "MAC地址不能为空",
	DEVICE_MAC_NO_FOUND_ERROR:        "请在管理平台注册设备",
	DEVICE_MAC_ERROR:                 "查询MAC信息错误",

	ERROR_USER_LOGIN_ERROR:           "用户名或者密码错误",
	ERROR_USER_NOT_EXIST:              "用户不存在",

	OLD_PASSWORD_ERROR:               "旧密码错误",

	CREATE_DIR_ERROR:                 "创建文件夹失败",


	WRITE_FILE_ERROR:                  "写入文件失败" ,
	READ_FILE_ERROR:                   "读取文件失败" ,
	InsertDataInfoError:               "写入data_info失败",
	GET_FILE_SIZE_ERROR:			   "获取文件大小失败",
	OPEN_FILE_ERROR:				   "打开文件失败",
	ILLEGAL_FILE_FORMAT:			   "文件格式不合法",
	CREATE_ZIP_ERROR:					"创建压缩文件失败",
	GET_FILE_HEADER_ERROR:				"获取文件头信息失败",
	WRITE_ZIP_ERROR:					"写入压缩文件失败",
	SAVE_FILE_ERROR:					"保存文件失败",
	FILE_SIZE_TOO_LARGE:				"文件大小过大",
	FILE_COMPRESS_ERROR:				"文件压缩失败",


}

func GetMsg(code string) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[SUCCESS]
}

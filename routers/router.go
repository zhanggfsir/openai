package routers

import (
	"openai-backend/controllers"

	"github.com/astaxie/beego"
)

func init() {
	filter()
	registerRouter()
	loginRouter()
	authRouter()
}

func registerRouter() {
	beego.Router("/", &controllers.MainController{})
	// auth get
	beego.Router("/api/:ver/gateway/redis/authcache", &controllers.AuthController{})
	//beego.Router("/api/:ver/gateway/logrecord", &controllers.AuthController{},"Post:LogRecord")

	//  monitorChart
	// API总调用量 调用情况统计		ApiTotalCalledPie
	beego.Router("/api/:ver/overView", &controllers.LogApiQuotasStatisticController{},"Get:ApiTotalCalledPie")
	// api日调用量 调用量日统计
	beego.Router("/api/:ver/apiChart", &controllers.LogApiDailyStatisticController{},"Get:ApiDailyCalledLine")

	// API用量统计 API调用总量
	beego.Router("/api/:ver/apiCalledTotal", &controllers.LogApiDailyStatisticController{},"Get:ApiCalledTotalList")

	// 应用列表 应用调用总量
	beego.Router("/api/:ver/appCalledTotal", &controllers.LogAppDailyStatisticController{},"Get:AppDailyCalledTotalList")

	//  调用总量趋势统计
	beego.Router("/api/:ver/appCalledTotalTrend", &controllers.LogHourlyQuotasStatisticController{},"Get:AppCalledTotalTrend")

	// 监控报表
	beego.Router("/api/:ver/monitorChart", &controllers.LogHourlyQuotasStatisticController{},"Get:AppApiMonitorChart")

	//  api 每日累计调用
	beego.Router("/api/:ver/apiCalledTotal/countLine", &controllers.LogApiDailyHistoryAddUpController{},"Get:ApiDailyHistoryAddUp")

	// 每个category 调用总量 成功 失败
	beego.Router("/api/:ver/apiCalledTotal/byCategory", &controllers.LogApiQuotasStatisticController{},"Get:CategoryPieChart")


	//  api被app调用情况 单个api调用情况（柱状图）
	beego.Router("/api/:ver/apiCalledTotal/apiCalledByAppHistogram", &controllers.LogDailyStatisticController{},"Get:ApiCalledByAppHistogram")

	//  app调用api情况 单个app调用情况（折线图）
	beego.Router("/api/:ver/appCallApiLineChart", &controllers.LogDailyStatisticController{},"Get:AppCallApiLineChart")


	//beego.Router("/api/:ver/categoryPieChart", &controllers.MonitorChartController{},"Get:categoryPieChart")


	beego.Router("/api/:ver/category", &controllers.AbilityCategoryController{})

	beego.Router("/api/:ver/abilityCategory", &controllers.AbilityCategoryController{})
	//能力API列表
	beego.Router("/api/:ver/ability", &controllers.AbilityController{})
	beego.Router("/api/:ver/api", &controllers.ApiController{})
	beego.Router("/api/:ver/apiList", &controllers.ApiController{}, "Get:GetApiList")

	//APP管理
	beego.Router("/api/:ver/app", &controllers.AppController{})
	beego.Router("/api/:ver/appList", &controllers.AppController{}, "Get:GetAppList")
	beego.Router("/api/:ver/appModification", &controllers.AppController{},"Post:Modify")
}

func loginRouter() {
	beego.Router("/api/:ver/captcha/:captchaId", &controllers.SysController{}, "Get:Captcha")
	beego.Router("/api/:ver/login", &controllers.SysController{}, "Post:Login")
	beego.Router("/api/:ver/logout", &controllers.SysController{}, "Post:Logout")
	beego.Router("/api/:ver/register", &controllers.SysController{}, "Post:Register")
	beego.Router("/api/:ver/retrievePassword", &controllers.SysController{}, "Post:RetrievePassword")

	beego.Router("/api/:ver/sms", &controllers.SysController{}, "Get:SmsCaptcha")
	beego.Router("/api/:ver/token/refresh", &controllers.SysController{}, "Get:RefreshToken")
	beego.Router("/api/:ver/thirdParty/token", &controllers.SysController{}, "Get:ThirdPartyToken")
}

func authRouter() {
	beego.Router("/api/:ver/role/authority", &controllers.RoleController{}, "Get:RoleAuthList")
	beego.Router("/api/:ver/role", &controllers.RoleController{})
	beego.Router("/api/:ver/role/list", &controllers.RoleController{}, "Get:GetList")
	//beego.Router("/api/:ver/role/policy/list", &controllers.RoleController{}, "Get:GetApiRoleAuthList")
	//beego.Router("/api/:ver/role/policy", &controllers.RoleController{}, "Put:UpdateApiRoleAuthList")

	beego.Router("/api/:ver/childUser", &controllers.ChildUserController{})
	beego.Router("/api/:ver/childUser/list", &controllers.ChildUserController{}, "Get:GetList")
	beego.Router("/api/:ver/childUser/role", &controllers.ChildUserController{}, "Put:ChildRoleManager")

	beego.Router("/api/:ver/policy", &controllers.PolicyController{})
	beego.Router("/api/:ver/policy/list", &controllers.PolicyController{}, "Get:GetList")

	beego.Router("/api/:ver/router", &controllers.RouterController{})
	beego.Router("/api/:ver/router/list", &controllers.RouterController{}, "Get:RouterList")
	beego.Router("/api/:ver/role/api/list", &controllers.SysApiController{}, "Get:GetRoleApiList")
	beego.Router("/api/:ver/role/api/updateCasbin", &controllers.SysApiController{},"Put:UpdateCasbin")
}

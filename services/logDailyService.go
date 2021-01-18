package services

import (
	"github.com/astaxie/beego/logs"
	"openai-backend/models"
	"strconv"
	"time"
)

func GetApiCalledByAppHistogram(queryFilter map[string]interface{}, domain string) (map[string]interface{}, error) {

	var jsonData = make(map[string]interface{})
	var labelList = make([]string, 0)
	var successfulList = make([]uint64, 0)
	var failedList = make([]uint64, 0)
	var sumList = make([]uint64, 0)

	histogramList, err := models.NewLogDailyQuotasStatisticModel().GetApiCalledByAppHistogramList(queryFilter, domain)
	if err != nil {
		logs.Error("get monitor chart error", err)
		return nil, err
	}


	for _, histogram := range histogramList {
		labelList = append(labelList, histogram.AppName)
		successfulList = append(successfulList, histogram.Successful)
		failedList = append(failedList, histogram.Failed)
		sumList = append(sumList, histogram.Sum)
	}

	// 1.得到所有label
	//labelList, err := models.NewHourlyQuotasStatisticModel().GetLabel(queryFilter,domain)
	//if err != nil {
	//	logs.Error("get monitor chart error", err)
	//	return nil, err
	//}

	jsonData["label"] = labelList

	// 2.得到所有data [内]
	//   2.1 得到 调用成功
	//var allData = make(map[string]interface{})
	var allData = make([]interface{}, 0)

	var success = make(map[string]interface{})
	//calledSuccess, err := models.NewHourlyQuotasStatisticModel().GetCalledSuccess(queryFilter,domain)
	//if err != nil {
	//	logs.Error("get monitor chart error", err)
	//	return nil, err
	//}

	//var jsonData = make(map[string]interface{})

	success["value"] = successfulList
	success["key"] = "调用成功"
	success["type"] = "bar"
	allData = append(allData, success)
	//allData["total"] = Count
	//allData["calledSuccessData"] = success

	//   2.2 得到 调用失败
	var failed = make(map[string]interface{})
	//calledFailed, err := models.NewHourlyQuotasStatisticModel().GetCalledFailed(queryFilter,domain)
	//if err != nil {
	//	logs.Error("get monitor chart error", err)
	//	return nil, err
	//}

	//var jsonData = make(map[string]interface{})
	failed["value"] = failedList
	failed["key"] = "调用失败"
	failed["type"] = "bar"
	//successAndFaiedData["total"] = Count
	allData = append(allData, failed)
	//allData["calledFailedData"] = failed

	//   2.3 得到 调用总数
	// GetCalledSuccessAndFailed
	var successAndFailed = make(map[string]interface{})

	//calledSuccessAndFailed, err := models.NewHourlyQuotasStatisticModel().GetCalledSuccessAndFailed(queryFilter,domain)
	//if err != nil {
	//	logs.Error("get monitor chart error", err)
	//	return nil, err
	//}

	//var jsonData = make(map[string]interface{})
	successAndFailed["value"] = sumList
	successAndFailed["key"] = "调用总数"
	successAndFailed["type"] = "bar"

	allData = append(allData, successAndFailed)
	//allData["calledFailedAndSuccess"] = successAndFailed

	//allData["total"] = Count

	jsonData["data"] = allData

	return jsonData, nil
}

/*

api_id api_name    day 		successful failed  sum
  30   人脸识别 2020-01-05  		1  		 0   	1
 166   智能门禁 2020-01-08		0		 20		20


 api_id  api_name    	day   	 	sum
  30   	 人脸识别 	2020-01-05	   	 1
  166    智能门禁 	2020-01-08	 	 20


第一步
api_id api_name    day 		   sum
  30   人脸识别 2020-01-05   	1
  30   人脸识别 2020-01-06 	 	0
  30   人脸识别 2020-01-07   	0
  30   人脸识别 2020-01-08   	0
  30   智能门禁 2020-01-05   	0
  30   智能门禁 2020-01-06 	 	0
  30   智能门禁 2020-01-07   	0
  30   智能门禁 2020-01-08   	0
 166   智能门禁 2020-01-08		20




*/

func GetAppCallApiLineChart(queryFilter map[string]interface{}, domain string) (map[string]interface{}, error) {

	var jsonData = make(map[string]interface{})
	//var labelList = make([]string, 0)

	lineChartList, err := models.NewLogDailyQuotasStatisticModel().GetAppCallApiLineChartList(queryFilter, domain)
	if err != nil {
		logs.Error("get monitor chart error", err)
		return nil, err
	}

	lineViewList:=make([]*models.AppCallApiLineChart,0)
	dateList := getBetweenDates(queryFilter["start_time"].(string), queryFilter["end_time"].(string) )
	//dateList:=getBetweenDates()
	lineChartMap:=make(map[string]interface{})
	for _, lineChart := range lineChartList {
		lineChartMap[lineChart.Day+strconv.FormatUint(lineChart.ApiId, 10)]=lineChart
	}

	// 会有重复问题  将放入lineChartModel map(day,lineChartModel)
	for _, date := range dateList {
		for _, lineChartModel := range lineChartList {
			//if lineChartModel.ApiId!=170{
			//	continue
			//}

			if date==lineChartModel.Day{
				lineViewList=append(lineViewList,lineChartModel)
			}else{
				// 实现 inner join 虚拟表，剔除重复元素的过程
				if lineChartMap[date+strconv.FormatUint(lineChartModel.ApiId, 10)]!=nil{  // map中已经存在，不再统计。重复不再计算 // lineChartModel.Day
					continue
				}

				appCallApiLineChart:=&models.AppCallApiLineChart{}
				appCallApiLineChart.ApiName=lineChartModel.ApiName
				appCallApiLineChart.ApiId=lineChartModel.ApiId
				appCallApiLineChart.Day=date
				appCallApiLineChart.Successful=0
				appCallApiLineChart.Failed=0
				appCallApiLineChart.Sum=0
				//logs.Info("******",appCallApiLineChart)
				lineChartMap[date+strconv.FormatUint(lineChartModel.ApiId, 10)]=appCallApiLineChart	// 再次加入map
				lineViewList=append(lineViewList,appCallApiLineChart)
			}
		}
	}



	// 会有重复问题  将放入lineChartModel map(day,lineChartModel)
	//for _, date := range dateList {
	//	for _, lineChartModel := range lineChartList {
	//		if date==lineChartModel.Day{
	//			lineViewList=append(lineViewList,lineChartModel)
	//			//logs.Info(date,"   -->",lineChartModel)
	//
	//		}else{
	//			logs.Info("****************************")
	//			logs.Info()
	//			logs.Info("****************************")
	//			appCallApiLineChart:=&models.AppCallApiLineChart{}
	//			appCallApiLineChart.ApiName=lineChartModel.ApiName
	//			appCallApiLineChart.ApiId=lineChartModel.ApiId
	//			appCallApiLineChart.Day=date
	//			appCallApiLineChart.Successful=0
	//			appCallApiLineChart.Failed=0
	//			appCallApiLineChart.Sum=0
	//			//logs.Info("******",appCallApiLineChart)
	//			lineViewList=append(lineViewList,appCallApiLineChart)
	//		}
	//	}
	//}



	apiMap := make(map[string][]uint64)

	for _, lineChart := range lineViewList {
		//for date
		//labelList = append(labelList, lineChart.Day)

		// k not exist then create a list ，then list add ；else list add
		if apiMap[lineChart.ApiName] == nil {
			apiMap[lineChart.ApiName] = make([]uint64, 0)
		}
		apiMap[lineChart.ApiName] = append(apiMap[lineChart.ApiName], lineChart.Sum)

	}

	// 1.得到所有label
	jsonData["label"] = dateList

	// 2 得到 dataList
	var dataList = make([]interface{}, 0)
	for apiName, apiSumList := range apiMap {
		var m = make(map[string]interface{})
		m["key"] = apiName
		m["type"] = "line"
		m["value"] = apiSumList
		dataList = append(dataList, m)
	}

	jsonData["data"] = dataList

	return jsonData, nil
}

// GetBetweenDates 根据开始日期和结束日期计算出时间段内所有日期
// 参数为日期格式，如：2020-01-01
// dateList:=getBetweenDates("2020-01-01","2020-01-10")
func getBetweenDates(startDate, endDate string) []string {
	d := make([]string, 0)
	//d := []string{}
	timeFormatTpl := "2006-01-02 15:04:05"
	if len(timeFormatTpl) != len(startDate) {
		timeFormatTpl = timeFormatTpl[0:len(startDate)]
	}
	date, err := time.Parse(timeFormatTpl, startDate)
	if err != nil {
		// 时间解析，异常
		return d
	}
	date2, err := time.Parse(timeFormatTpl, endDate)
	if err != nil {
		// 时间解析，异常
		return d
	}
	if date2.Before(date) {
		// 如果结束时间小于开始时间，异常
		return d
	}
	// 输出日期格式固定
	timeFormatTpl = "2006-01-02"
	date2Str := date2.Format(timeFormatTpl)
	d = append(d, date.Format(timeFormatTpl))
	for {
		date = date.AddDate(0, 0, 1)
		dateStr := date.Format(timeFormatTpl)
		d = append(d, dateStr)
		if dateStr == date2Str {
			break
		}
	}
	return d
}


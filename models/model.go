package models

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/olivere/elastic/v7"
	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v1/s3"
	"log"
	"math"
	"os"
	"strings"
	"time"
)

type PublicModel struct {
	CreateTime time.Time  `orm:"column(created_time);type(datetime);auto_now_add" json:"created_time"`
	ModifyTime time.Time  `orm:"column(updated_time);type(datetime);auto_now" json:"updated_time"`
}

type DelRequest struct {
	Id []int `json:"id"`
}

var (
	UserTypeHost = "1"
	UserTypeChild = "2"
)

var (
	EsClient *elastic.Client
	CephConn *s3.S3

)

func InitElastic() {
	//es 配置
	var err error
	isOpenElastic := beego.AppConfig.DefaultBool("elastic", false)
	if !isOpenElastic {
		return
	}
	host := beego.AppConfig.DefaultString("elastic_host", "")
	hosts := strings.Split(host, ",")
	elastic.SetSniff(false)
	errorLog := log.New(os.Stdout, "APP", log.LstdFlags)
	user := beego.AppConfig.DefaultString("elastic_user", "elastic")
	password := beego.AppConfig.DefaultString("elastic_password", "@#&qazXSW")
	EsClient, err = elastic.NewClient(elastic.SetErrorLog(errorLog), elastic.SetURL(hosts...),elastic.SetBasicAuth(user,password))
	if err != nil {
		logs.Error("Elastic缓存初始化失败")
		//panic(err)
	}
	_, err = EsClient.ElasticsearchVersion(host)
	if err != nil {
		logs.Error("Elastic缓存初始化失败")
		panic(err)
	}
	logs.Info("Elastic缓存初始化完成")
}



//https://www.cnblogs.com/xdao/p/elasticsearch_golang.html
func elasticFilter(tableName string,
	queryFilter map[string]interface{}) *elastic.SearchService{
	result := EsClient.Search(tableName)
	if len(queryFilter) >0 {
		boolQ := elastic.NewBoolQuery()
		for key, value := range queryFilter {
			if key == "temperature" {
				spliceEsQuery(boolQ, key, value)
			} else if key == "snap_shot_time" {
				spliceEsQuery(boolQ, key, value)
			} else if key == "domain" {
				boolQ.Must(elastic.NewTermQuery(key, value))
			} else {
				//boolQ.Must(elastic.NewMatchQuery(key, value))
				//elastic.NewMatchPhraseQuery()
				boolQ.Must(elastic.NewMatchPhraseQuery(key, value))
			}
		}
		result = result.Query(boolQ)
	}
	return result
}

//拼接es查询语句
func spliceEsQuery(boolQ *elastic.BoolQuery, key string,
	value interface{}) (*elastic.BoolQuery, error) {
	var mapValue map[string]interface{}
	switch value.(type) {
	case map[string]interface{}:
		mapValue = value.(map[string]interface{})
	default:
		return nil, errors.New("参数错误")
	}
	if len(mapValue) == 2 {
		return boolQ.Filter(elastic.NewRangeQuery(key).Gte(mapValue["start"]),
			elastic.NewRangeQuery(key).Lte(mapValue["end"])), nil
	} else if _, startOk := mapValue["start"]; startOk {
		return boolQ.Filter(elastic.NewRangeQuery(
			key).Gte(mapValue["start"])), nil
	} else if _, endOk := mapValue["end"]; endOk {
		return boolQ.Filter(elastic.NewRangeQuery(
			key).Lte(mapValue["end"])), nil
	} else {
		return nil, errors.New("参数错误")
	}
}



//获得es插入的索引
func esIndex(tableName string) string {
	return tableName + "_" + time.Now().Format("200601")
}


func getOffset(pageSize int,  pageIndex int, valueNum int64) int {
	var offset int
	totalPage := math.Ceil(float64(valueNum)/float64(pageSize)) // 总共需要多少页 向上取整
	if pageIndex <= int(totalPage) {
		offset = (pageIndex - 1) * pageSize  // 第pageIndex页的起始 offset
	} else {
		offset = (int(totalPage) - 1) * pageSize
	}
	return offset
}

//func RegisterUserPolicy(domain string) error {
//	adminPolicy := make(map[string][]string)
//	list, err := NewPolicyModel().FetchAllInfo()
//	if err != nil || len(list) == 0 {
//		logs.Error("fetch all policy", utils.Fields{
//			"err": err,
//			"count": len(list),
//			"domain": domain,
//		})
//
//		return errors.New("fetch all policy error")
//	}
//
//	for _, policyInfo := range list{
//		adminPolicy[policyInfo.Policy] = casbinUtil.DefaultActionList
//	}
//
//	return AddPolicyFromController(domain, casbinUtil.RoleAdmin, adminPolicy);
//}



func InitCeph() {
	// 初始化ceph信息
	auth := aws.Auth{  //"XAZJTJ4KK7DI56L1MTAG", "vPqoGKoSJEpUuwWdF7fxZijQpjWNV6CqB5rrqTWX"
		AccessKey: beego.AppConfig.DefaultString("access_key", ""),
		SecretKey: beego.AppConfig.DefaultString("secret_key", ""),
	}
	region := aws.Region{
		Name: 		 beego.AppConfig.DefaultString("ceph_region_name", ""),
		EC2Endpoint: beego.AppConfig.DefaultString("ec2_end_point", ""),
		S3Endpoint:  beego.AppConfig.DefaultString("s3_end_point", ""),
		S3BucketEndpoint:"",
		S3LocationConstraint:false, // 没有区域限制
		S3LowercaseBucket:false, // bucket没有大小写限制
		Sign:aws.SignV2,
	}
	// 创建🔓s3类型连接
	CephConn = s3.New(auth, region)
}




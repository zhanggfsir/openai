package casbinUtil

/*
@Author:
@Time: 2020-07-31 10:02
@Description: casbin
*/

type CasbinRule struct {
	Id    int    // 自增主键
	PType string // Policy Type - 用于区分 policy和 group(role)
	V0    string // subject
	V1    string // object
	V2    string // action
	V3    string // 这个和下面的字段无用，仅预留位置，如果你的不是
	V4    string // sub, obj, act的话才会用到
	V5    string // 如 sub, obj, act, suf就会用到 V3
}

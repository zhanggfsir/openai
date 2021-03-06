package casbinUtil

import (
	"github.com/astaxie/beego/orm"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"runtime"
)

/*
@Author:
@Time: 2020-07-31 10:02
@Description: casbin
*/

var Enforcer *casbin.Enforcer

type Adapter struct {
	o orm.Ormer
}

func RegisterCasbin() {
	a := &Adapter{}
	a.o = orm.NewOrm()
	runtime.SetFinalizer(a, finalizer)
	// Enforcer初始化 - 即传入Adapter对象
	var err error
	Enforcer, err = casbin.NewEnforcer("conf/casbin.conf", a)
	if err != nil {
		panic(err)
	}
	// Enforcer读取Policy
	err = Enforcer.LoadPolicy()
	if err != nil {
		panic(err)
	}
}

// LoadPolicy loads policy from database.
func (a *Adapter) LoadPolicy(model model.Model) error {
	var lines []CasbinRule
	_, err := a.o.QueryTable("casbin_rule").All(&lines)
	if err != nil {
		return err
	}

	for _, line := range lines {
		loadPolicyLine(line, model)
	}

	return nil
}

func loadPolicyLine(line CasbinRule, model model.Model) {
	lineText := line.PType
	if line.V0 != "" {
		lineText += ", " + line.V0
	}
	if line.V1 != "" {
		lineText += ", " + line.V1
	}
	if line.V2 != "" {
		lineText += ", " + line.V2
	}
	if line.V3 != "" {
		lineText += ", " + line.V3
	}
	if line.V4 != "" {
		lineText += ", " + line.V4
	}
	if line.V5 != "" {
		lineText += ", " + line.V5
	}

	persist.LoadPolicyLine(lineText, model)
}

func savePolicyLine(ptype string, rule []string) CasbinRule {
	line := CasbinRule{}

	line.PType = ptype
	if len(rule) > 0 {
		line.V0 = rule[0]
	}
	if len(rule) > 1 {
		line.V1 = rule[1]
	}
	if len(rule) > 2 {
		line.V2 = rule[2]
	}
	if len(rule) > 3 {
		line.V3 = rule[3]
	}
	if len(rule) > 4 {
		line.V4 = rule[4]
	}
	if len(rule) > 5 {
		line.V5 = rule[5]
	}

	return line
}

// SavePolicy saves policy to database.
func (a *Adapter) SavePolicy(model model.Model) error {
	var lines []CasbinRule

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			lines = append(lines, line)
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			lines = append(lines, line)
		}
	}

	_, err := a.o.InsertMulti(len(lines), lines)
	return err
}

// AddPolicy adds a policy rule to the storage.
func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	line := savePolicyLine(ptype, rule)
	_, err := a.o.Insert(&line)
	return err
}

// RemovePolicy removes a policy rule from the storage.
func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	line := savePolicyLine(ptype, rule)
	_, err := a.o.Delete(&line, "p_type", "v0", "v1", "v2", "v3", "v4", "v5")
	return err
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	line := CasbinRule{}

	line.PType = ptype
	filter := []string{}
	filter = append(filter, "p_type")
	if fieldIndex <= 0 && 0 < fieldIndex+len(fieldValues) {
		line.V0 = fieldValues[0-fieldIndex]
		filter = append(filter, "v0")
	}
	if fieldIndex <= 1 && 1 < fieldIndex+len(fieldValues) {
		line.V1 = fieldValues[1-fieldIndex]
		filter = append(filter, "v1")
	}
	if fieldIndex <= 2 && 2 < fieldIndex+len(fieldValues) {
		line.V2 = fieldValues[2-fieldIndex]
		filter = append(filter, "v2")
	}
	if fieldIndex <= 3 && 3 < fieldIndex+len(fieldValues) {
		line.V3 = fieldValues[3-fieldIndex]
		filter = append(filter, "v3")
	}
	if fieldIndex <= 4 && 4 < fieldIndex+len(fieldValues) {
		line.V4 = fieldValues[4-fieldIndex]
		filter = append(filter, "v4")
	}
	if fieldIndex <= 5 && 5 < fieldIndex+len(fieldValues) {
		line.V5 = fieldValues[5-fieldIndex]
		filter = append(filter, "v5")
	}
	_, err := a.o.Delete(&line, filter...)
	return err
}

func (a *Adapter) createTable() {
	err := orm.RunSyncdb("casbin", false, true)
	if err != nil {
		panic(err)
	}
}

func (a *Adapter) dropTable() {
	err := orm.RunSyncdb("casbin", true, true)
	if err != nil {
		panic(err)
	}
}

// finalizer is the destructor for Adapter.
func finalizer(_ *Adapter) {}

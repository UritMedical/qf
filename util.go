// Package qf
// @Description: 框架的一些通用方法
package qf

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
)

//
// buildTableName
//  @Description: 根据结构体，生成对应的数据库表名
//  @param model 结构体
//  @return string 然后表名，规则：包名_结构体名，如果包名和结构体名一致时，则只返回结构体名
//
func buildTableName(model interface{}) string {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	pName := strings.ToLower(filepath.Base(t.PkgPath()))
	bName := strings.ToLower(t.Name())
	tName := fmt.Sprintf("%s_%s", pName, bName)
	if pName == bName {
		tName = pName
	}
	return tName
}

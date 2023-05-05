package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/UritMedical/qf/util/qreflect"
	"reflect"
	"strings"
)

//
// ToMap
//  @Description: 将包含内容的结构体转为字典结构
//  @param model 内容结构
//  @return map[string]interface{} 字典
//
func ToMap(model interface{}) map[string]interface{} {
	ref := qreflect.New(model)
	return ref.ToMapExpandAll()
}

//
// ToMaps
//  @Description: 将内容列表转为字典列表
//  @param list 内容结构列表
//  @return []map[string]interface{} 字典列表
//
func ToMaps(list interface{}) []map[string]interface{} {
	values := reflect.ValueOf(list)
	if values.Kind() != reflect.Slice {
		panic(fmt.Errorf("list must be slice"))
	}
	finals := make([]map[string]interface{}, values.Len())
	for i := 0; i < values.Len(); i++ {
		finals[i] = ToMap(values.Index(i).Interface())
	}
	return finals
}

//
// Bind
//  @Description: 将source结构体中的数据绑定到target结构体中
//  @param targetPtr vm结构体, 必须是指针
//  @param source 原始结构体
//  @return error
//
func Bind(targetPtr interface{}, source interface{}) error {
	r := qreflect.New(source)
	return SetModel(targetPtr, r.ToMapExpandAll())
}

//
// SetModel
//  @Description: 修改结构体内的字段值
//  @param objectPtr
//  @param value
//
func SetModel(objectPtr interface{}, value map[string]interface{}) error {
	if objectPtr == nil {
		return errors.New("the object cannot be empty")
	}
	ref := qreflect.New(objectPtr)
	// 必须为指针
	if ref.IsPtr() == false {
		return errors.New("the object must be pointer")
	}

	// 修改外部值
	if value != nil {
		e := ref.SetAny(value)
		if e != nil {
			return e
		}
	}
	// 修改Info
	return setInfo(ref, value)
}

//
// SetList
//  @Description: 修改列表
//  @param objectPtr
//  @param values
//  @return error
//
func SetList(objectPtr interface{}, values []map[string]interface{}) error {
	if objectPtr == nil {
		return errors.New("the objectPtr cannot be empty")
	}
	ref := qreflect.New(objectPtr)
	// 必须为指针
	if ref.IsSlice() == false {
		return errors.New("the objectPtr must be slice")
	}

	// 修改外部值
	if values != nil {
		e := ref.SetAny(values)
		if e != nil {
			return e
		}
	}
	// 修改Info
	objs := ref.InterfaceArray()
	for i, obj := range objs {
		e := setInfo(qreflect.New(obj), values[i])
		if e != nil {
			return e
		}
	}
	ref.Clear()
	return ref.SetAny(objs)
}

func setInfo(ref *qreflect.Reflect, value map[string]interface{}) error {
	all := ref.ToMap()

	// 复制一份
	temp := map[string]interface{}{}
	for k, v := range value {
		temp[k] = v
	}

	// 转摘要
	if field, ok := temp["SummaryFields"]; ok && field != "" {
		e := ref.Set("Summary", fields(field, all["Summary"], all, &temp))
		if e != nil {
			return e
		}
	}
	// 转信息
	if field, ok := temp["InfoFields"]; ok && field != "" {
		e := ref.Set("Info", fields(field, all["Info"], all, &temp))
		if e != nil {
			return e
		}
		return nil
	}

	// 将剩余的全部写入到Info中
	if info, ok := all["Info"]; ok {
		mp := map[string]interface{}{}
		_ = json.Unmarshal([]byte(info.(string)), &mp)
		for k, v := range temp {
			if k == "SummaryFields" || k == "InfoFields" {
				continue
			}
			if _, ok := all[k]; ok == false {
				mp[k] = v
			}
		}
		mj, _ := json.Marshal(mp)
		e := ref.Set("Info", string(mj))
		if e != nil {
			return e
		}
	}
	return nil
}

func fields(field interface{}, source interface{}, all map[string]interface{}, values *map[string]interface{}) string {
	if field == nil || field.(string) == "" {
		return ""
	}
	// 获取原始数据并转为字典
	mp := map[string]interface{}{}
	if source != nil {
		_ = json.Unmarshal([]byte(source.(string)), &mp)
	}
	// 获取需要的值
	temp := *values
	for _, name := range strings.Split(field.(string), ",") {
		if _, ok := all[name]; ok == false {
			if _, ok2 := temp[name]; ok2 {
				mp[name] = temp[name]
				delete(temp, name)
			}
		}
	}
	values = &temp
	// 返回
	if len(mp) == 0 {
		return ""
	}
	mj, _ := json.Marshal(mp)
	return string(mj)
}

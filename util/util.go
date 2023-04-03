package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/UritMedical/qf/util/qreflect"
	"reflect"
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
		e := ref.Set(value)
		if e != nil {
			return e
		}
	}
	// 修改FullInfo
	return setFullInfo(ref, value)
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
		e := ref.Set(values)
		if e != nil {
			return e
		}
	}
	// 修改FullInfo
	objs := ref.InterfaceArray()
	for i, obj := range objs {
		e := setFullInfo(qreflect.New(obj), values[i])
		if e != nil {
			return e
		}
	}
	ref.Clear()
	return ref.Set(objs)
}

func setFullInfo(ref *qreflect.Reflect, value map[string]interface{}) error {
	all := ref.ToMap()
	if info, ok := all["FullInfo"]; ok {
		str := info.(string)
		mp := map[string]interface{}{}
		err := json.Unmarshal([]byte(str), &mp)
		if err == nil || str == "" {
			for k, v := range value {
				if _, ok := all[k]; ok == false {
					mp[k] = v
				}
			}
			mj, _ := json.Marshal(mp)
			e := ref.Set(map[string]interface{}{"FullInfo": string(mj)})
			if e != nil {
				return e
			}
		}
	}
	return nil
}

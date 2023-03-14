package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/UritMedical/qf/util/reflectex"
	"reflect"
	"time"
)

//
// ToMap
//  @Description: 将包含内容的结构体转为字典结构
//  @param model 内容结构
//  @return map[string]interface{} 字典
//
func ToMap(model interface{}) map[string]interface{} {
	// 先转一次json
	tj, _ := json.Marshal(model)
	// 然后在反转到内容对象
	cnt := struct {
		Id       uint64
		LastTime time.Time
		FullInfo string
	}{}
	_ = json.Unmarshal(tj, &cnt)

	// 生成字典
	final := join(cnt.FullInfo, model)
	// 补齐字段的值
	final["Id"] = cnt.Id

	return final
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
//  @Description: 修改结构体内的方法
//  @param objectPtr
//  @param value
//
func SetModel(objectPtr interface{}, value map[string]interface{}) error {
	if objectPtr == nil {
		return errors.New("the object cannot be empty")
	}
	// 必须为指针
	if reflectex.IsPtr(objectPtr) == false {
		return errors.New("the object must be pointer")
	}
	mp := reflectex.StructToMap(objectPtr)
	for k, v := range value {
		if _, ok := mp[k]; ok {
			mp[k] = v
		}
	}
	if info, ok := mp["FullInfo"]; ok {
		imap := map[string]interface{}{}
		err := json.Unmarshal([]byte(info.(string)), &imap)
		if err == nil {
			for k, v := range value {
				if _, ok := mp[k]; ok == false {
					imap[k] = v
				}
			}
			mj, _ := json.Marshal(imap)
			mp["FullInfo"] = string(mj)
		}
	}
	return nil
}

// 将完整内容Json和对应的实体，合并为一个字典对象
func join(info string, model interface{}) map[string]interface{} {
	data := map[string]interface{}{}

	// 将内容的信息写入到字典中
	_ = json.Unmarshal([]byte(info), &data)

	// 反射对象，并将其他字段附加到字典
	value := reflect.ValueOf(model)
	if value.Kind() == reflect.Map {
		for _, v := range value.MapKeys() {
			data[v.String()] = value.MapIndex(v).Interface()
		}
	} else {
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}
		for i := 0; i < value.NumField(); i++ {
			field := value.Field(i)
			// 通过原始内容
			if field.Kind() == reflect.Struct && field.Type().Name() == "BaseModel" {
				continue
			}
			tag := value.Type().Field(i).Tag.Get("json")
			if tag != "-" {
				data[value.Type().Field(i).Name] = field.Interface()
			}
		}
	}
	return data
}

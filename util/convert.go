package util

import (
	"encoding/json"
	"fmt"
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

func SetModel(model interface{}, value map[string]interface{}) {

}

// 将完整内容Json和对应的实体，合并为一个字典对象
func join(info string, model interface{}) map[string]interface{} {
	data := map[string]interface{}{}

	// 将内容的信息写入到字典中
	_ = json.Unmarshal([]byte(info), &data)

	// 反射对象，并将其他字段附加到字典
	value := reflect.ValueOf(model)
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

	return data
}

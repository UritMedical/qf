package util

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/UritMedical/qf/util/qreflect"
	"math/rand"
	"reflect"
	"time"
)

//ConvertToMD5 转换成MD5加密
func ConvertToMD5(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

//
// DiffIntSet
//  @Description: 计算a数组元素不在b数组之中的所有元素
//  @param a
//  @param b
//  @return []uint64
//
func DiffIntSet(a []uint64, b []uint64) []uint64 {
	c := make([]uint64, 0)
	temp := map[uint64]struct{}{}
	//把b所有的值作为key存入temp
	for _, val := range b {
		if _, ok := temp[val]; !ok {
			temp[val] = struct{}{}
		}
	}
	//如果a中的值作为key在temp中找不到，说明它不在b中
	for _, val := range a {
		if _, ok := temp[val]; !ok {
			c = append(c, val)
		}
	}
	return c
}

// RandomString
//  @Description: 生成随机字符串
//  @param length
//  @return string
//
func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	result := make([]rune, length)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}

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
		LastTime uint64
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
	// 修改FullInfo值
	all := ref.ToMap()
	if info, ok := all["FullInfo"]; ok {
		mp := map[string]interface{}{}
		err := json.Unmarshal([]byte(info.(string)), &mp)
		if err == nil {
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

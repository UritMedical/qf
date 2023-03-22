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

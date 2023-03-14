package reflectex

import (
	"reflect"
)

//
// IsPtr
//  @Description: 判断对象是否是指针
//  @param obj
//  @return bool
//
func IsPtr(obj interface{}) bool {
	return reflect.ValueOf(obj).Kind() == reflect.Ptr
}

//
// IsStruct
//  @Description: 判断对象是否是结构体
//  @param obj
//  @return bool
//
func IsStruct(obj interface{}) bool {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.Kind() == reflect.Struct
}

//
// IsSlice
//  @Description: 判断对象是否是切片
//  @param obj
//  @return bool
//
func IsSlice(obj interface{}) bool {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.Kind() == reflect.Slice
}

//
// StructToMap
//  @Description: 将对象结构体转为字典
//  @param obj
//  @return map[string]string
//
func StructToMap(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	output := map[string]interface{}{}
	if t.Kind() == reflect.Map {
		for _, m := range v.MapKeys() {
			output[m.String()] = v.MapIndex(m).Interface()
		}
	} else {
		recursionStructToMap(output, t, v)
	}
	return output
}

// 递归将对象结构体转为字典
func recursionStructToMap(output map[string]interface{}, t reflect.Type, v reflect.Value) {
	if v.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
		// 如果是空列表，则默认扩充一条用于反射结构体内部的字段
		if t.Kind() == reflect.Slice && v.Len() == 0 {
			v = reflect.MakeSlice(t, 1, 1)
			v = v.Index(0)
			t = v.Type()
		}
	}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Struct {
			if field.String() == "<time.Time Value>" {
				output[t.Field(i).Name] = field.Interface()
			} else {
				recursionStructToMap(output, field.Type(), field)
			}
		} else if field.CanInterface() {
			output[t.Field(i).Name] = field.Interface()
		}
	}
}

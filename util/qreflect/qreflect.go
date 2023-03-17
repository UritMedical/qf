package qreflect

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Reflect struct {
	t  reflect.Type
	v  reflect.Value
	kv map[string]interface{}
}

//
// New
//  @Description: 创建反射实例
//  @param object 任意对象
//  @return Reflect
//
func New(object interface{}) *Reflect {
	r := &Reflect{
		t: reflect.TypeOf(object),
		v: reflect.ValueOf(object),
	}
	return r
}

//
// IsPtr
//  @Description: 判断是否是指针
//  @receiver r
//  @return bool
//
func (r *Reflect) IsPtr() bool {
	return r.v.Kind() == reflect.Ptr
}

//
// ToMap
//  @Description: 转为字典
//  @return map[string]interface{}
//
func (r *Reflect) ToMap() map[string]interface{} {
	if r.kv == nil {
		r.getMap(r.t, r.v)
	}
	return r.kv
}

//
// Set
//  @Description: 将结构体或字典写入到对象中，支持结构体和切片
//  @param map[string]interface{}
//  @return error
//
func (r *Reflect) Set(values ...interface{}) error {
	if values == nil || len(values) == 0 {
		return errors.New("the value 's length must be greater than 0")
	}
	// 如果对象不是指针，则无法执行
	if r.v.Kind() != reflect.Ptr {
		return errors.New("the obj 's kind must be ptr")
	}
	// 如果是结构体
	kind := r.t.Elem().Kind()
	if kind == reflect.Struct {
		for _, value := range values {
			e := r.setStruct(value)
			if e != nil {
				return e
			}
		}
		return nil
	}
	// 如果是切片
	if kind == reflect.Slice {
		// 先计算长度，如果全部长度不一致，则失败
		vLen := 0
		vValue := make([]reflect.Value, 0)
		for _, value := range values {
			t := reflect.TypeOf(value)
			v := reflect.ValueOf(value)
			// 如果参数不是切片，则失败
			if t.Kind() != reflect.Slice {
				return errors.New("the value 's type must be slice")
			}
			if vLen != 0 && vLen != v.Len() {
				return errors.New("the length of the parameter is inconsistent")
			}
			vLen = v.Len()
			vValue = append(vValue, v)
		}
		slice := reflect.MakeSlice(r.t.Elem(), vLen, vLen)
		for _, value := range vValue {
			e := r.setSlice(slice, value)
			if e != nil {
				return e
			}
		}
		r.v.Elem().Set(reflect.AppendSlice(r.v.Elem(), slice))
		return nil
	}
	return errors.New("the obj 's kind must be struct or slice")
}

func (r *Reflect) setStruct(sub interface{}) error {
	subT := reflect.TypeOf(sub)
	subV := reflect.ValueOf(sub)
	// 如果入参是切片，则只取第一条
	if subT.Kind() == reflect.Slice {
		if subV.Len() == 0 {
			return errors.New("the value 's length must be greater than 0")
		}
		subT = subV.Index(0).Type()
		subV = subV.Index(0)
	}
	v := r.v.Elem()
	return r.set(v, subT, subV)
}

func (r *Reflect) setSlice(v reflect.Value, sub reflect.Value) error {
	for i := 0; i < sub.Len(); i++ {
		vv := v.Index(i)
		st := sub.Index(i).Type()
		sv := sub.Index(i)
		e := r.set(vv, st, sv)
		if e != nil {
			return e
		}
	}
	return nil
}

func (r *Reflect) set(v reflect.Value, subT reflect.Type, subV reflect.Value) error {
	if subT.Kind() == reflect.Map {
		for _, m := range subV.MapKeys() {
			f := v.FieldByName(m.String())
			if f.CanSet() {
				vv, err := r.convert(f.Type().String(), subV.MapIndex(m).Interface())
				if err == nil {
					f.Set(reflect.ValueOf(vv))
				}
			}
		}
	} else {
		if subT.Kind() == reflect.Ptr {
			subT = subT.Elem()
			subV = subV.Elem()
		}
		for i := 0; i < subV.NumField(); i++ {
			f := v.FieldByName(subT.Field(i).Name)
			if f.CanSet() {
				tp := strings.ToLower(f.Type().Kind().String())
				vv, err := r.convert(tp, subV.Field(i).Interface())
				if err == nil {
					// 如果是自定义类型包原生类型，则特殊处理
					if tp != f.Type().String() {
						switch tp {
						case "uint", "uint8", "uint32", "uint64":
							u, e := strconv.ParseUint(fmt.Sprintf("%v", vv), 10, 64)
							if e == nil {
								f.SetUint(u)
							}
						case "int", "int8", "int32", "int64":
							u, e := strconv.ParseInt(fmt.Sprintf("%v", vv), 10, 64)
							if e == nil {
								f.SetInt(u)
							}
						}
					} else {
						f.Set(reflect.ValueOf(vv))
					}
				}
			}
		}
	}
	return nil
}

// 转为字典
func (r *Reflect) getMap(t reflect.Type, v reflect.Value) {
	if r.kv == nil {
		r.kv = map[string]interface{}{}
	}
	if t.Kind() == reflect.Map {
		for _, m := range v.MapKeys() {
			r.kv[m.String()] = v.MapIndex(m).Interface()
		}
	} else {
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
					r.kv[t.Field(i).Name] = field.Interface()
				} else {
					r.getMap(field.Type(), field)
				}
			} else if field.CanInterface() {
				r.kv[t.Field(i).Name] = field.Interface()
			}
		}
	}
}

// 类型转换
func (r *Reflect) convert(typeName string, value interface{}) (interface{}, error) {
	str := fmt.Sprintf("%v", value)
	switch typeName {
	case "string":
		return str, nil
	case "int":
		return strconv.Atoi(str)
	case "int64":
		return strconv.ParseInt(str, 10, 64)
	case "int32":
		v, e := strconv.ParseInt(str, 10, 32)
		return int32(v), e
	case "int16":
		v, e := strconv.ParseInt(str, 10, 16)
		return int16(v), e
	case "int8":
		v, e := strconv.ParseInt(str, 10, 8)
		return int8(v), e
	case "uint":
		v, e := strconv.ParseUint(str, 10, 64)
		return uint(v), e
	case "uint64":
		return strconv.ParseUint(str, 10, 64)
	case "uint32":
		v, e := strconv.ParseUint(str, 10, 32)
		return uint32(v), e
	case "uint16":
		v, e := strconv.ParseUint(str, 10, 16)
		return uint16(v), e
	case "uint8":
		v, e := strconv.ParseUint(str, 10, 8)
		return uint8(v), e
	case "bool":
		return strconv.ParseBool(str)
	case "float64":
		return strconv.ParseFloat(str, 64)
	case "float32":
		v, e := strconv.ParseFloat(str, 32)
		return float32(v), e
	case "time.Time":
		return value, nil
	case "time.Duration":
		v, e := strconv.Atoi(str)
		return time.Duration(v) * time.Millisecond, e
	case "*string":
		return &str, nil
	case "*int":
		v, e := strconv.Atoi(str)
		return &v, e
	case "*int64":
		v, e := strconv.ParseInt(str, 10, 64)
		return &v, e
	case "*int32":
		v, e := strconv.ParseInt(str, 10, 32)
		o := int32(v)
		return &o, e
	case "*int16":
		v, e := strconv.ParseInt(str, 10, 16)
		o := int16(v)
		return &o, e
	case "*int8":
		v, e := strconv.ParseInt(str, 10, 8)
		o := int8(v)
		return &o, e
	case "*uint":
		v, e := strconv.ParseUint(str, 10, 64)
		o := uint(v)
		return &o, e
	case "*uint64":
		v, e := strconv.ParseUint(str, 10, 64)
		return &v, e
	case "*uint32":
		v, e := strconv.ParseUint(str, 10, 32)
		o := uint32(v)
		return &o, e
	case "*uint16":
		v, e := strconv.ParseUint(str, 10, 16)
		o := uint16(v)
		return &o, e
	case "*uint8":
		v, e := strconv.ParseUint(str, 10, 8)
		o := uint8(v)
		return &o, e
	case "*bool":
		v, e := strconv.ParseBool(str)
		return &v, e
	case "*float64":
		v, e := strconv.ParseFloat(str, 64)
		return &v, e
	case "*float32":
		v, e := strconv.ParseFloat(str, 32)
		o := float32(v)
		return &o, e
	case "*time.Time":
		return &value, nil
	case "*time.Duration":
		v, e := strconv.Atoi(str)
		o := time.Duration(v) * time.Millisecond
		return &o, e
	}
	return nil, errors.New("Not support Type " + typeName)
}
package qf

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// EKind 行为类别
type EKind string

var (
	EKindSave     EKind = "Save"     // 新增或修改
	EKindDelete   EKind = "Delete"   // 删除
	EKindGetModel EKind = "GetModel" // 获取单条
	EKindGetList  EKind = "GetList"  // 获取多条
)

//
// HttpMethod
//  @Description: 返回Http方式名称
//  @return string
//
func (kind EKind) HttpMethod() string {
	if kind == EKindSave {
		return "POST"
	}
	if kind == EKindDelete {
		return "DELETE"
	}
	return "GET"
}

//
// Content
//  @Description: 基础内容实体对象
//
type Content struct {
	Id     uint64    `gorm:"primaryKey"` // 唯一号
	Delete byte      `gorm:"index"`      // 是否删除 0否 1是
	Time   time.Time `gorm:"index"`      // 操作时间
	Info   string    // 完整内容信息
}

//
// Context
//  @Description: Api上下文参数
//
type Context struct {
	Time     time.Time // 操作时间
	UserId   uint      // 操作用户账号
	UserName string    // 操作用户名字

	jsonValue   map[string]interface{}
	stringValue string
}

func (ctx *Context) Bind(object interface{}) error {
	if object == nil {
		return errors.New("object is empty")
	}

	// 反转json
	err := json.Unmarshal([]byte(ctx.stringValue), object)
	if err != nil {
		return err
	}
	// 反射对象
	value := reflect.ValueOf(object)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	// 判断类型
	switch value.Type().Kind() {
	case reflect.Struct: // 结构体
		source := map[string]interface{}{}
		if json.Unmarshal([]byte(ctx.stringValue), &source) == nil {
			nj, _ := json.Marshal(ctx.build(source))
			_ = json.Unmarshal(nj, object)
		}
	case reflect.Slice: // 列表
		source := make([]map[string]interface{}, 0)
		if json.Unmarshal([]byte(ctx.stringValue), &source) == nil {
			cnt := make([]Content, 0)
			for i := 0; i < len(cnt); i++ {
				cnt = append(cnt, ctx.build(source[i]))
			}
			nj, _ := json.Marshal(cnt)
			_ = json.Unmarshal(nj, object)
		}
	}
	return nil
}

func (ctx *Context) build(source map[string]interface{}) Content {
	nid := 0
	if id, ok := source["Id"]; ok {
		v, e := strconv.Atoi(fmt.Sprintf("%v", id))
		if e == nil {
			nid = v
		}
	}
	if nid == 0 {
		// 通过平台获取ID
	}
	cj, _ := json.Marshal(source)
	return Content{
		Id:   uint64(nid),
		Time: ctx.Time,
		Info: string(cj),
	}
}

func (ctx *Context) GetJsonValue(propName string) string {
	obj := ctx.jsonValue[propName]
	if obj == nil {
		return ""
	}
	nj, _ := json.Marshal(obj)
	return string(nj)
}

func (ctx *Context) GetStringValue(propName string) string {
	obj := ctx.jsonValue[propName]
	if obj == nil {
		return ""
	}
	return fmt.Sprintf("%v", obj)
}

func (ctx *Context) GetUIntValue(propName string) uint64 {
	num, _ := strconv.Atoi(ctx.GetStringValue(propName))
	return uint64(num)
}

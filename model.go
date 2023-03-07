package qf

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// EApiKind 行为类别
type EApiKind string

var (
	EApiKindSave     EApiKind = "Save"     // 新增或修改
	EApiKindDelete   EApiKind = "Delete"   // 删除
	EApiKindGetModel EApiKind = "GetModel" // 获取单条
	EApiKindGetList  EApiKind = "GetList"  // 获取多条
)

//
// HttpMethod
//  @Description: 返回Http方式名称
//  @return string
//
func (kind EApiKind) HttpMethod() string {
	if kind == EApiKindSave {
		return "POST"
	}
	if kind == EApiKindDelete {
		return "DELETE"
	}
	return "GET"
}

//
// BaseModel
//  @Description: 基础实体对象
//
type BaseModel struct {
	Id       uint64    `gorm:"primaryKey"` // 唯一号
	LastTime time.Time `gorm:"index"`      // 最后操作时间时间
	FullInfo string    // 内容
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

func (ctx *Context) Bind(objectPtr interface{}) error {
	if objectPtr == nil {
		return errors.New("the object cannot be empty")
	}
	// 必须为指针
	value := reflect.ValueOf(objectPtr)
	if value.Kind() != reflect.Ptr {
		return errors.New("the object must be pointer")
	}

	// 反转json
	err := json.Unmarshal([]byte(ctx.stringValue), objectPtr)
	if err != nil {
		return err
	}
	value = value.Elem()
	// 判断类型
	switch value.Type().Kind() {
	case reflect.Struct: // 结构体
		source := map[string]interface{}{}
		if json.Unmarshal([]byte(ctx.stringValue), &source) == nil {
			nj, _ := json.Marshal(ctx.build(source))
			_ = json.Unmarshal(nj, objectPtr)
		}
	case reflect.Slice: // 列表
		source := make([]map[string]interface{}, 0)
		if json.Unmarshal([]byte(ctx.stringValue), &source) == nil {
			cnt := make([]BaseModel, 0)
			for i := 0; i < len(cnt); i++ {
				cnt = append(cnt, ctx.build(source[i]))
			}
			nj, _ := json.Marshal(cnt)
			_ = json.Unmarshal(nj, objectPtr)
		}
	}
	return nil
}

func (ctx *Context) GetMap() map[string]interface{} {
	return nil
}

func (ctx *Context) BindWithMap(objectPtr interface{}, joins ...interface{}) error {
	return nil
}

func (ctx *Context) build(source map[string]interface{}) BaseModel {
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
	return BaseModel{
		Id:       uint64(nid),
		LastTime: ctx.Time,
		FullInfo: string(cj),
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

func (ctx *Context) GetId() uint64 {
	return ctx.GetUIntValue("Id")
}

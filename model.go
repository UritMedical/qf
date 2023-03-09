package qf

import (
	"encoding/json"
	"errors"
	"fmt"
	"qf/util/reflectex"
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
	UserId   uint64    // 操作用户账号
	UserName string    // 操作用户名字

	// input的原始内容字典
	inputValue  []map[string]interface{}
	inputSource string
	// id分配器
	idPer       uint
	idAllocator iIdAllocator
}

//
// Bind
//  @Description: 绑定到新实体对象
//  @param objectPtr 实体对象指针（必须为指针）
//  @param joins 需要附加的值（可以是结构体、字典）
//  @return error
//
func (ctx *Context) Bind(objectPtr interface{}, joins ...interface{}) error {
	if objectPtr == nil {
		return errors.New("the object cannot be empty")
	}
	// 必须为指针
	if reflectex.IsPtr(objectPtr) == false {
		return errors.New("the object must be pointer")
	}

	// 追加附加内容到字典
	for _, join := range joins {
		for k, v := range reflectex.StructToMap(join) {
			for _, vv := range ctx.inputValue {
				vv[k] = v
			}
		}
	}
	// 然后根据类型，将字典写入到对象或列表中
	table := buildTableName(objectPtr)
	cnt := make([]BaseModel, 0)
	for i := 0; i < len(ctx.inputValue); i++ {
		c := ctx.build(ctx.inputValue[i], reflectex.StructToMap(objectPtr))
		// 如果Id为0，则自动分配信息Id
		if c.Id == 0 {
			c.Id = ctx.idAllocator.Next(table)
		}
		cnt = append(cnt, c)
	}
	if reflectex.IsStruct(objectPtr) {
		// 先将提交的input填充
		nj, _ := json.Marshal(ctx.inputValue[0])
		_ = json.Unmarshal(nj, objectPtr)
		// 再将重新组织的内容填充
		nj, _ = json.Marshal(cnt[0])
		_ = json.Unmarshal(nj, objectPtr)
	} else if reflectex.IsSlice(objectPtr) {
		// 同上
		nj, _ := json.Marshal(ctx.inputValue)
		_ = json.Unmarshal(nj, objectPtr)
		nj, _ = json.Marshal(cnt)
		_ = json.Unmarshal(nj, objectPtr)
	}
	return nil
}

func (ctx *Context) build(source map[string]interface{}, exclude map[string]interface{}) BaseModel {
	nid := uint64(0)
	if id, ok := source["Id"]; ok {
		v, e := strconv.Atoi(fmt.Sprintf("%v", id))
		if e == nil {
			nid = uint64(v)
		}
	}
	// 从完整的原始input中，去掉实体对象中已经存在的
	finals := map[string]interface{}{}
	for k, v := range source {
		if _, ok := exclude[k]; ok == false {
			finals[k] = v
		}
	}
	cj, _ := json.Marshal(finals)
	return BaseModel{
		Id:       nid,
		LastTime: ctx.Time,
		FullInfo: string(cj),
	}
}

func (ctx *Context) GetJsonValue(propName string) string {
	if len(ctx.inputValue) == 0 {
		return ""
	}
	nj, _ := json.Marshal(ctx.inputValue[0][propName])
	return string(nj)
}

func (ctx *Context) GetStringValue(propName string) string {
	if len(ctx.inputValue) == 0 {
		return ""
	}
	return fmt.Sprintf("%v", ctx.inputValue[0][propName])
}

func (ctx *Context) GetUIntValue(propName string) uint64 {
	num, _ := strconv.Atoi(ctx.GetStringValue(propName))
	return uint64(num)
}

func (ctx *Context) GetId() uint64 {
	return ctx.GetUIntValue("Id")
}

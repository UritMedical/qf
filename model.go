package qf

import "time"

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
	ID   uint      `gorm:"primarykey"` // 唯一号
	Time time.Time `gorm:"index"`      // 操作时间
	Info string    // 完整内容信息
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

func (ctx *Context) BindModel(model interface{}) error {
	return nil
}

func (ctx *Context) GetStringValue(propName string) string {
	return ""
}

func (ctx *Context) GetUIntValue(propName string) uint {
	return 0
}

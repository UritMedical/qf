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

// Content 基础内容实体
type Content struct {
	ID   uint   `gorm:"primarykey"` // 唯一号
	Info string // 完整内容信息
}

// Context 上下文
type Context struct {
	Time     time.Time // 操作时间
	UserId   uint      // 操作用户账号
	UserName string    // 操作用户名字
}

func (ctx Context) BindModel(model interface{}) {

}

func (ctx Context) ToContent(model interface{}) Content {
	return Content{}
}

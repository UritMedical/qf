package qf

import (
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
	Id       uint64    `gorm:"primaryKey;autoIncrement:false"` // 唯一号
	LastTime time.Time `gorm:"index"`                          // 最后操作时间时间
	FullInfo string    // 内容
}

//
// LoginUser
//  @Description: 登陆用户信息
//
type LoginUser struct {
	UserId      uint64 // 登陆用户唯一号
	UserName    string // 登陆用户名字
	LoginId     string // 登陆用户账号
	Departments map[uint64]struct {
		Name string
	} // 所属部门列表
	token string // 登陆的token信息
	roles map[uint64]struct {
		Name string
	} // 角色列表
}

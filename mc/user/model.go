package user

import "qf"

// User
// @Description: 用户
//
type User struct {
	qf.Content
	LoginId  string `gorm:"uniqueIndex"` //登录Id
	Password string //密码
}

// UserRole
// @Description: 用户角色关系
//
type UserRole struct {
	Id     uint
	UserId uint `gorm:"index"`
	RoleId uint `gorm:"index"`
}

// Role
// @Description: 角色
//
type Role struct {
	Id   uint
	Name string // 角色名称
}

// RoleRights
// @Description: 角色权限组关系
//
type RoleRights struct {
	Id       uint
	RoleId   uint `gorm:"index"`
	RightsId uint `gorm:"index"`
}

// RightsGroup
// @Description: 权限组
//
type RightsGroup struct {
	Id   uint
	Name string //权限组名称
}

// RightsApi
// @Description: 权限组与API的关系
//
type RightsApi struct {
	Id       uint
	RightsId uint   `gorm:"index"`
	ApiId    string //API key
}

// Department
// @Description: 部门
// Info 里面包含子组织
type Department struct {
	Id       uint
	Name     string // 部门名称
	ParentId uint   `gorm:"index"` //父级部门Id
}

// DepartUser
// @Description: 用户组织关系表
//
type DepartUser struct {
	Id         uint
	OrganizeId uint `gorm:"index"`
	UserId     uint `gorm:"index"`
}

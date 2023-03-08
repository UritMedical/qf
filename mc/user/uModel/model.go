package uModel

import "qf"

// User
// @Description: 用户
//
type User struct {
    qf.BaseModel
    LoginId  string `gorm:"unique"` //登录Id
    Password string `json:"-"`      //密码
}

// UserRole
// @Description: 用户角色关系
//
type UserRole struct {
    qf.BaseModel
    UserId uint64 `gorm:"index"`
    RoleId uint64 `gorm:"index"`
}

// Role
// @Description: 角色
//
type Role struct {
    qf.BaseModel
    Name string `gorm:"unique"` // 角色名称
}

// RoleRights
// @Description: 角色权限组关系
//
type RoleRights struct {
    qf.BaseModel
    RoleId   uint64 `gorm:"index"`
    RightsId uint64 `gorm:"index"`
}

// RightsGroup
// @Description: 权限组
//
type RightsGroup struct {
    qf.BaseModel
    Name string `gorm:"unique"` //权限组名称
}

// RightsApi
// @Description: 权限组与API的关系
//
type RightsApi struct {
    qf.BaseModel
    RightsId uint64 `gorm:"index"`
    ApiId    string //API key
}

// Department
// @Description: 部门
// Info 里面包含子组织
type Department struct {
    qf.BaseModel
    Name     string `gorm:"unique"` // 部门名称
    ParentId uint64 `gorm:"index"`  //父级部门Id
}

// DepartUser
// @Description: 用户组织关系表
//
type DepartUser struct {
    qf.BaseModel
    DepartId uint64 `gorm:"index"`
    UserId   uint64 `gorm:"index"`
}
